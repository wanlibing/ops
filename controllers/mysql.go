package controllers

import (
	"github.com/astaxie/beego/orm"
	"msql/models"
	"strconv"
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
)

//默认首页登录入口控制器
type MysqlController struct {
	BaseController
}


//删除commit_status表数据
//重构使用于多种情况删除，user用户表
//不带主键不能删除？？
func (c *MysqlController) DeleteIdSubmitData(){

	token,_  := c.ParseHeaderToken()
	claim := token.Claims
	fmt.Println(claim)


	Msql := orm.NewOrm()
	Msql.Using("default")
	var resultNum,status string


	table := c.GetString("table")
	if table == "submit_status" {
		sqlid,_ := strconv.Atoi(c.GetString("id"))
		submitinfo := models.Submit_status{Id:sqlid}

		//Msql.QueryTable("Submit_status").Filter("id",id).One(&submitinfo,"Id","Dbname","Sqlstatment","Submittime","Approver","ApprovalStatus","Operator","OperatorStatus")
		if _,err := Msql.Delete(&submitinfo) ; err != nil {
			resultNum = err.Error()
			status = "failed"
		} else {
			resultNum = fmt.Sprintf("已删除行数：","1")
			status = "successed"
		}
	} else if table == "user"  {
		username := c.GetString("username")
		fmt.Println("user name isssssssssss",username)
		userinfo := models.User{Name:username}
		fmt.Println(userinfo)
		//Msql.Delete(&userinfo)   这种默认会采用主键删除
		Msql.QueryTable("User").Filter("Name", username).Delete()
	}
	resultNum = "fuck you"
	resultRow := struct {
		Status string
		Num string
	}{status,resultNum}
	c.Data["json"] = &resultRow
	c.ServeJSON()
}


//新增sql审计任务，表插入数据
func (c *MysqlController) InsertSubmitData(){
	fmt.Println("debug here now 111111")
	token,_  := c.ParseHeaderToken()
	claim := token.Claims.(jwt.MapClaims)
	fmt.Println(claim)

	fmt.Println("debug curuser  now is ...",claim["username"])


	dbname := c.GetString("dbname")
	checkername := c.GetString("checkername")
	sql := c.GetString("sqlstatement")
	submittime := time.Now().Format("2006-01-02 15:04:05")

	Msql := orm.NewOrm()
	Msql.Using("default")


	//执行人是谁，如果有多个运维是不是应该发给整个运维组？
	insertSqlJob  := models.Submit_status{Dbname:dbname,Approver:checkername,Sqlstatment:sql,Submittime:submittime,Submitter:claim["username"].(string)}  //类型强制转换
	Msql.Insert(&insertSqlJob)

	result := struct {
		Status string
	}{"ok"}
	c.Data["json"] = &result
	c.ServeJSON()
}

//审批催促
func (c *MysqlController) UrgedSubmitData(){
	pusher := c.GetString("pusher")
	sqlid,_ := strconv.Atoi(c.GetString("sqlid"))
	receiver := c.GetString("reciver")

	Msql := orm.NewOrm()
	Msql.Using("default")
	//执行人是谁，如果有多个运维是不是应该发给整个运维组？
	//insertSqlJob  := models.Submit_status{Dbname:dbname,Approver:checkername,Sqlstatment:sql,Submittime:submittime}
	insertUeged := models.SqlPushRelationship{Pusher:pusher,SqlId:sqlid,Receiver:receiver}
	Msql.Insert(&insertUeged)

	//将推送信息存入隧道
	tmpUrgedMsg := urgedmsg{Pusher:pusher,Sqlid:sqlid,Receiver:receiver}
	urgedmeschan <- tmpUrgedMsg
	fmt.Println("save urged msg sucesssssssssed")
	result := struct {
		Status string
	}{"ok"}
	c.Data["json"] = &result
	c.ServeJSON()
}

//审批者接受审批提醒,如何实现？？
//应该是用redis，那种读取消费，或者用chan实现
//采用basecontroller，然后增加一个chan，已不行，容易堵塞，前端循环读的或容易堵住
//前端轮训取服务端信息，会比较耗费资源，简单粗暴，但是a用户会读取b用户的chan
//轮训取推送信息的ajax放在哪里？
func (c *MysqlController) ReceiveUrgedSubmitData(){
	tmpUrgedMsg :=  urgedmsg{}
	fmt.Println("begingetmsg")
	select {
	case tmpUrgedMsg = <- urgedmeschan: //拿到锁
		fmt.Println("get ugred  msg")
	case <-time.After(5 * time.Second): //超时5s
		fmt.Println("dong get ugred msg ")
	}

	c.Data["json"] = &tmpUrgedMsg
	c.ServeJSON()
}


//审核或者执行状态改变
//重构，能更新所有表
func (c *MysqlController) UpdateSubmitData(){
	Msql := orm.NewOrm()
	Msql.Using("default")
	sqlid,_ := strconv.Atoi(c.GetString("sqlid"))

	if c.GetString("checkstatus") != ""{
		checkstatus,_ := strconv.Atoi(c.GetString("checkstatus"))
		//将ApprovalStatus字段值设置为1，审核成功
		Msql.QueryTable("Submit_status").Filter("Id",sqlid).Update(orm.Params{"ApprovalStatus":checkstatus})
	} else if c.GetString("execstatus") != ""{
		execstatus,_ := strconv.Atoi(c.GetString("execstatus"))
		//将ApprovalStatus字段值设置为1，审核成功
		Msql.QueryTable("Submit_status").Filter("Id",sqlid).Update(orm.Params{"OperatorStatus":execstatus})
	} else {
		fmt.Println("update job status nothing ")
	}
	result := struct {
		Status string
	}{"ok"}
	c.Data["json"] = &result
	c.ServeJSON()
}

//获取SQL审计用户数据
func (c *MysqlController) GetSqlUserManager(){
	var sqluser []models.User
	Msql := orm.NewOrm()
	Msql.Using("default")
	//id倒叙排序，前端展示最新提交信息
	Msql.QueryTable("User").All(&sqluser,"Name","Department","Roleid")
	c.Data["json"] = &sqluser
	c.ServeJSON()
}