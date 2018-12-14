package routers

import (
	"msql/controllers"
	"github.com/astaxie/beego"
)



func init() {
	
	//beego.SetStaticPath("/static","staic")
    beego.Router("/login", &controllers.MainController{})
    beego.Router("/login/handler",&controllers.LoginController{},"POST:LoginStatus")
	beego.Router("/home/submitter",&controllers.SubmitterController{})
	beego.Router("/home/submitter/access",&controllers.SubmitterController{},"GET:SubmitAccess")
	beego.Router("/home/submitter/processing",&controllers.SubmitterController{},"GET:SubmitProcessing")
	beego.Router("/home/checker",&controllers.CheckerController{})
	beego.Router("/home/ops",&controllers.OpsController{})
	//运维看板
	beego.Router("/home/ops/baseenv",&controllers.OpsenvController{})
    //运维管理平台
	beego.Router("/ops/manager/monitor",&controllers.OpsController{},"GET:OpsManagerMonitor")
	beego.Router("/ops/manager/deployment",&controllers.OpsController{},"GET:OpsManagerDeployment")

    //监控系统
	beego.Router("/api/proxy/query",&controllers.OpsController{},"GET:GetQueryPrometheusData")
    //发布系统
	beego.Router("/api/deploy/additem",&controllers.OpsController{},"POST:AddDeploymentItem")
	beego.Router("/api/deploy/getitemplan",&controllers.OpsController{},"GET:GetItemPlan")

	beego.Router("/mysql/delete",&controllers.MysqlController{},"GET:DeleteIdSubmitData")
	beego.Router("/mysql/insert",&controllers.MysqlController{},"POST:InsertSubmitData")


	beego.Router("/mysql/update",&controllers.MysqlController{},"POST:UpdateSubmitData")
	beego.Router("/mysql/urged",&controllers.MysqlController{},"POST:UrgedSubmitData")
	beego.Router("/mysql/urged",&controllers.MysqlController{},"GET:ReceiveUrgedSubmitData")
    //获取SQL审计用户数据
	beego.Router("/mysql/getsqlusermanager",&controllers.MysqlController{},"GET:GetSqlUserManager")

    //用户管理页面
	beego.Router("/user/sqlmanagement",&controllers.UserController{},"GET:SqlUserManagement")

    //用户管理
	beego.Router("/user/sql/edit",&controllers.UserController{},"GET:SqlUserUpdate")
	beego.Router("/user/sql/add",&controllers.UserController{},"GET:SqlUserAdd")

    //kubernetes 管理页面
	beego.Router("/home/k8s/nodepage",&controllers.KubernetesController{},"GET:GetNodePage")
	beego.Router("/api/k8s/namespaces",&controllers.KubernetesController{},"GET:GetNamespace")
	beego.Router("/home/k8s/namespace/:ns.*",&controllers.KubernetesController{},"GET:GetNamespaceInformation")


}
