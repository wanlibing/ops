package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"msql/models"
	"fmt"
	 "net/url"
	"github.com/astaxie/beego/httplib"
	"strconv"
	"reflect"
	"time"
	"strings"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//默认首页登录入口控制器
type OpsController struct {
	BaseController
}

type prometheusData struct {
	Status string
	Data struct{
		Result []struct{
			//Value []interface{}
			Value []interface{}
		}
	}
}

type prometheusRangeData struct {
	Status string
	Data struct{
		Result []struct{
			//Value []interface{}
			Values [][]interface{}
		}
	}
}

type responseProxyPrometheusMonitorData struct {
	DiskUsagedPrecent string
	MemoryUsagedPrecent string
	NodeSystemLoad []string
	TimestrampX []string
	NetworkReceive []string
	NetworkTransmit []string
	DiskRead []string
	DiskWrite []string
}

func (c *OpsController) Get() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))
	c.getAllSubmitData()
	//管理员页面获取k8s namespaces  treeview

	c.Data["username"] = "运维管理员"
	c.Data["thisName"]="运维管理员"
	c.TplName = "role/ops/ops.html"
	c.Layout = "common/layout_home.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["treeviewhtml"] = "role/ops/treeview.html"
}

//采用排序来控制在页面上的展示最新的任务，并用js实现是否能撤销操作
func (c *OpsController) getAllSubmitData(){
	var submitinfo []models.Submit_status
	Msql := orm.NewOrm()
	Msql.Using("default")
	//id倒叙排序，前端展示最新提交信息
	Msql.QueryTable("Submit_status").All(&submitinfo,"Id","Dbname","Sqlstatment","Submittime","Approver","ApprovalStatus","Operator","OperatorStatus")
	c.Data["AllSubmitinfo"] = submitinfo

}


func (c *OpsController) OpsManagerMonitor() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))
	c.TplName = "opsmanager/monitor.html"
	c.Layout = "common/layout_home.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["treeviewhtml"] = "role/ops/treeview.html"
}

func (c *OpsController) GetQueryPrometheusData(){
	url := beego.AppConfig.String("prometheus_server")
	uri := beego.AppConfig.String("prometheus_query")
	uriRange := beego.AppConfig.String("prometheus_query_range")
	monitorIp := c.GetString("ip")
	var queryParamsLoad []string
	var cpuLoad []string

	queryload1mParam := fmt.Sprintf("node_load1{instance=~%s}",monitorIp)
	queryParamsLoad = append(queryParamsLoad,queryload1mParam)
	queryload5mParam := fmt.Sprintf("node_load5{instance=~%s}",monitorIp)
	queryParamsLoad = append(queryParamsLoad,queryload5mParam)
	queryload15mParam := fmt.Sprintf("node_load15{instance=~%s}",monitorIp)
	queryParamsLoad = append(queryParamsLoad,queryload15mParam)

	for _,query := range queryParamsLoad{
		httpUrl := urlEncode(url + uri + query)
		cpuloadreq := httplib.Get(httpUrl)
		cpuloadresult := prometheusData{}
		cpuloadreq.ToJSON(&cpuloadresult)
		cpuLoad = append(cpuLoad,cpuloadresult.Data.Result[0].Value[1].(string))
	}
	fmt.Println("cpuload is",cpuLoad)


	//磁盘，内存使用百分比

	var queryParamsMemoryAndFilesystem []string
	//totalmemory,availmemory,total device sda ,avail device sda
	var memoryAndFilesystem []int

	queryMemoryTotalParam := fmt.Sprintf("node_memory_MemTotal_bytes{instance=~%s}",monitorIp)
	queryParamsMemoryAndFilesystem = append(queryParamsMemoryAndFilesystem,queryMemoryTotalParam)
	queryMemoryAvailParam := fmt.Sprintf("node_memory_MemAvailable_bytes{instance=~%s}",monitorIp)
	queryParamsMemoryAndFilesystem = append(queryParamsMemoryAndFilesystem,queryMemoryAvailParam)
	queryFilesystemTotalParam := fmt.Sprintf("node_filesystem_size_bytes{instance=~%s,device=~'.*sda.*'}",monitorIp)
	queryParamsMemoryAndFilesystem = append(queryParamsMemoryAndFilesystem,queryFilesystemTotalParam)
	queryFilesystemAvailParam := fmt.Sprintf("node_filesystem_avail_bytes{instance=~%s,device=~'.*sda.*'}",monitorIp)
	queryParamsMemoryAndFilesystem = append(queryParamsMemoryAndFilesystem,queryFilesystemAvailParam)

	for _,query := range queryParamsMemoryAndFilesystem{
		httpUrl := urlEncode(url + uri + query)
		memoryAndFilesystemReq := httplib.Get(httpUrl)
		memoryAndFilesystemResult := prometheusData{}
		memoryAndFilesystemReq.ToJSON(&memoryAndFilesystemResult)
		//cpuLoad = append(cpuLoad,cpuloadresult.Data.Result[0].Value[1].(string))
		tmpInt,_ := strconv.Atoi(memoryAndFilesystemResult.Data.Result[0].Value[1].(string))
		memoryAndFilesystem = append(memoryAndFilesystem,tmpInt )
	}

	fmt.Println(memoryAndFilesystem)

	//diskUsageRate,_ := fmt.Printf("%.2f",1 - float32(memoryAndFilesystem[3]) / float32(memoryAndFilesystem[2]))
	//memoryUsageRate,_ := fmt.Printf("%.2f",1 - float32(memoryAndFilesystem[1]) / float32(memoryAndFilesystem[0]))
	diskUsageRate := fmt.Sprintf("%.2f",1 - float32(memoryAndFilesystem[3]) / float32(memoryAndFilesystem[2]))
	memoryUsageRate := fmt.Sprintf("%.2f",1 - float32(memoryAndFilesystem[1]) / float32(memoryAndFilesystem[0]))

	fmt.Println(diskUsageRate, "and type is ",reflect.TypeOf(diskUsageRate))
	fmt.Println(memoryUsageRate)


	//device io and network io
	end := time.Now().Unix()
	//end := 1535095249
	start := end - 3600    //最近一小時

	var diskIoAndNetworkIo []string
	var timeX int64    //时间戳
	var hourMinutesX string
	var xAxisHourMinutes []string
	var yAxisIoBytes [][]string //network receive transmit disk read write (sda)
	var xAxisTimestamp  []string

	//取最近一小时，没隔5分钟统计一次
	queryRangeNetworkReceiveParam := fmt.Sprintf("rate (node_network_receive_bytes_total{device=~'ens33.*',instance=~%s}[1m])&start=%d&end=%d&step=300",monitorIp,start,end)
	diskIoAndNetworkIo = append(diskIoAndNetworkIo,queryRangeNetworkReceiveParam)
	queryRangeNetworkTransmitParam := fmt.Sprintf("rate (node_network_transmit_bytes_total{device=~'ens33.*',instance=~%s}[1m])&start=%d&end=%d&step=300",monitorIp,start,end)
	diskIoAndNetworkIo = append(diskIoAndNetworkIo,queryRangeNetworkTransmitParam)
	queryRangeDiskReadParam := fmt.Sprintf("rate (node_disk_read_bytes_total{device=~'sda',instance=~%s}[1m])&start=%d&end=%d&step=300",monitorIp,start,end)
	diskIoAndNetworkIo = append(diskIoAndNetworkIo,queryRangeDiskReadParam)
	queryRangeDiskWriteParam := fmt.Sprintf("rate (node_disk_written_bytes_total{device=~'sda',instance=~%s}[1m])&start=%d&end=%d&step=300",monitorIp,start,end)
	diskIoAndNetworkIo = append(diskIoAndNetworkIo,queryRangeDiskWriteParam)
	//fmt.Println("diskIoAndNetworkIo is" , diskIoAndNetworkIo)
	//fmt.Println("~~~~~~~~~~~")
	for _,queryRange := range diskIoAndNetworkIo {
		httpUrlRange := urlEncode(url + uriRange + queryRange)
		diskAndNetworkReq := httplib.Get(httpUrlRange)
		diskAndNetworkResult := prometheusRangeData{}
		diskAndNetworkReq.ToJSON(&diskAndNetworkResult)
		networkIoData :=  diskAndNetworkResult.Data.Result[0].Values
		var yAxisIo []string

		xAxisTimestamp = []string{}
		xAxisHourMinutes = []string{}
		//fmt.Println("debug befor",xAxisTimestamp)
		for _,v := range networkIoData{
			timex := v[0]
			timeX = int64(timex.(float64))
			hourMinutesX = time.Unix(timeX,0).Format("15:04 ")
			netwokBytesTmp := v[1]
			netwokBytes , _ := strconv.ParseFloat(netwokBytesTmp.(string),64)
			netwokBytesY := fmt.Sprintf("%.2f",netwokBytes)
			//fmt.Println(uint64(timex.(float64)))
			xAxisTimestamp = append(xAxisTimestamp,hourMinutesX)
			yAxisIo = append(yAxisIo,netwokBytesY)

			//fmt.Println(timeX)
			//fmt.Println("netwokBytes is ",netwokBytesY)
		}
		yAxisIoBytes = append(yAxisIoBytes,yAxisIo)
		//fmt.Println("xAxisTimestamp is" , xAxisTimestamp)
		xAxisHourMinutes = xAxisTimestamp
		//fmt.Println("debug after",xAxisHourMinutes)

	}

	//fmt.Println("x轴",xAxisHourMinutes)
	//fmt.Println("y轴",yAxisIoBytes)

	finalResultData := responseProxyPrometheusMonitorData{
		DiskUsagedPrecent: diskUsageRate,
		MemoryUsagedPrecent: memoryUsageRate,
		NodeSystemLoad: cpuLoad,
		TimestrampX: xAxisHourMinutes,
		NetworkReceive: yAxisIoBytes[0],
		NetworkTransmit: yAxisIoBytes[1],
		DiskRead: yAxisIoBytes[2],
		DiskWrite: yAxisIoBytes[3],
	}


	//c.Ctx.WriteString(cpuLoad[1])
	c.Data["json"] = &finalResultData
	c.ServeJSON()

}

func  urlEncode(u string) (urlparse string) {
	urlParse,_ := url.Parse(u)
	param := urlParse.Query()
	urlParse.RawQuery = param.Encode()
//	fmt.Println("urlparse is",urlParse)
	return urlParse.String()

}


//发布系统
//发布项目
func (c *OpsController) OpsManagerDeployment() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))
	c.GetDeploymentItem()
	c.TplName = "opsmanager/deployment.html"
	c.Layout = "common/layout_home.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["treeviewhtml"] = "role/ops/treeview.html"
}


type DeploymentText struct {
	//Id bson.ObjectId `bson:"_id"`
	Title string
	DateTime string
	Content string
	DeploymentStep []string
	SqlDbname []SqlPlan
}

type DeploymentTitle struct {
	//Id bson.ObjectId `bson:"_id"`
	Title string
	DateTime string
}

type SqlPlan map[string][]string

func (c *OpsController) GetDeploymentItem() {
	mongoServer := beego.AppConfig.String("mongodb_endpoint")
	mongoSession,err := mgo.Dial(mongoServer)

	if err != nil {
		fmt.Println("cannot conn to mongodb ",mongoServer)
		fmt.Println("this err is ",err)
		//mess := "connect to mongodb server failed"
		c.Ctx.WriteString("connect to mongodb server failed")
		//c.Ctx.Redirect(302,"/static/err/500.html")
		return
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)


	var job DeploymentText

	fmt.Println("this inline job us",job)
	getJobFromMongodb := mongoSession.DB(beego.AppConfig.String("mongodb_db")).C(beego.AppConfig.String("mongodb_table"))

	result := []DeploymentTitle{}
	err = getJobFromMongodb.Find(bson.M{}).All(&result)
	if err != nil {
		//log.Fatal(err)
	}
	c.Data["jobitem"] = result
}



func (c *OpsController) AddDeploymentItem() {
	onlineTitle	:= c.GetString("onlineTitle")
	onlineContent := c.GetString("onlineContent")
//	onlineItems	:= c.GetString("onlineItems")
	onlineItems	:= strings.Split(c.GetString("onlineItems"),",")
	onlineDbname := c.GetString("onlineDbname")
//	sql語句必須正确语法，比如必须有分号分隔符，否则将抛弃这些
	var onlineSqlStatement []string
	tmpsqlslice := strings.Split(c.GetString("onlineSqlStatement"),";")
	for _,tmpsql := range tmpsqlslice {
		if strings.TrimSpace(tmpsql) != "" {
			onlineSqlStatement = append(onlineSqlStatement,strings.TrimSpace(tmpsql))
		}
	}

	fmt.Println("onlineSqlStatement isssssssssss",onlineSqlStatement )

	//onlineSqlStatement := strings.Split(c.GetString("onlineSqlStatement"),";")


	mongoServer := beego.AppConfig.String("mongodb_endpoint")
	mongoSession,err := mgo.Dial(mongoServer)

	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	var sqlstatement SqlPlan
	if len(onlineSqlStatement) < 1 {
	} else {
		sqlstatement = SqlPlan{onlineDbname:onlineSqlStatement}
	}
	var job DeploymentText

	if sqlstatement == nil {
		job = DeploymentText{
			Title: onlineTitle,
			DateTime: time.Now().Format("2006-01-02 15:04"),
			Content: onlineContent,
			DeploymentStep: onlineItems,
		}
	} else {
		job = DeploymentText{
			Title: onlineTitle,
			DateTime: time.Now().Format("2006-01-02 15:04"),
			Content: onlineContent,
			DeploymentStep: onlineItems,
			SqlDbname: []SqlPlan{sqlstatement},
		}
	}
	fmt.Println("this inline job us",job)
	jobToMongodb := mongoSession.DB(beego.AppConfig.String("mongodb_db")).C(beego.AppConfig.String("mongodb_table"))

	err = jobToMongodb.Insert(&job)
	if err != nil {
		//log.Fatal(err)
	}

	result := struct {
		Status string
	}{"ok"}
	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *OpsController) GetItemPlan() {
	jobtime := c.GetString("lastjobtime")

	mongoServer := beego.AppConfig.String("mongodb_endpoint")
	mongoSession,err := mgo.Dial(mongoServer)

	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)



	getJobFromMongodb := mongoSession.DB(beego.AppConfig.String("mongodb_db")).C(beego.AppConfig.String("mongodb_table"))

	result := DeploymentText{}
	err = getJobFromMongodb.Find(bson.M{"datetime":jobtime}).One(&result)
	if err != nil {
		//log.Fatal(err)
	}
	c.Data["json"] = &result
	c.ServeJSON()
}


