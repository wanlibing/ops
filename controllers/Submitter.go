package controllers

import (
	"github.com/astaxie/beego/orm"
	"msql/models"
	"github.com/dgrijalva/jwt-go"
)

//默认首页登录入口控制器
type SubmitterController struct {
	//beego.Controller
	BaseController
}

func (c *SubmitterController) Get() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))
	//验证token
	token := c.ParseParamsToken()
	claim := token.Claims.(jwt.MapClaims)

	//fmt.Println("debug now is ...",claim["username"]) 取值
	c.Data["username"] = claim["username"]
	c.Data["thisName"]="大前端"


	c.getAllSubmitData()

	c.TplName = "role/submitter/submitter.html"
	c.Layout = "common/layout_home.html"
}

//采用排序来控制在页面上的展示最新的任务，并用js实现是否能撤销操作
func (c *SubmitterController) getAllSubmitData(){
	var submitinfo []models.Submit_status

	Msql := orm.NewOrm()
	Msql.Using("default")
	//id倒叙排序，前端展示最新提交信息
	Msql.QueryTable("Submit_status").Filter("Submitter",c.Data["username"]).All(&submitinfo,"Id","Dbname","Sqlstatment","Submittime","Approver","ApprovalStatus","Operator","OperatorStatus")
	c.Data["AllSubmitinfo"] = submitinfo

	c.Data["JobAllNum"] ,_ = Msql.QueryTable("Submit_status").Filter("Submitter",c.Data["username"]).Count()
	c.Data["JobNoCheckNum"] ,_ = Msql.QueryTable("Submit_status").Filter("Submitter",c.Data["username"]).Filter("ApprovalStatus",0).Count()
	c.Data["JobNoExecNum"] ,_ = Msql.QueryTable("Submit_status").Filter("Submitter",c.Data["username"]).Filter("OperatorStatus",0).Count()

}

func (c *SubmitterController) SubmitAccess() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))
	//验证token


	//fmt.Println("debug now is ...",claim["username"]) 取值
	c.Data["username"] = "wanlb"
	c.Data["thisName"]="大前端"


	c.getAccessSubmitData()
	//这里应该取到用户的userid再决定下面的定义
	c.TplName = "role/submitter/accessSubmitter.html"
	//c.Layout = "common/layout_home.html"
	c.Layout = "common/layout_home.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["treeviewhtml"] = "role/ops/treeview.html"
}
func (c *SubmitterController) getAccessSubmitData(){
	var submitinfo []models.Submit_status

	Msql := orm.NewOrm()
	Msql.Using("default")
	//id倒叙排序，前端展示最新提交信息
	Msql.QueryTable("Submit_status").Filter("Submitter",c.Data["username"]).Filter("ApprovalStatus",1).All(&submitinfo,"Id","Dbname","Sqlstatment","Submittime","Approver","ApprovalStatus","Operator","OperatorStatus")
	c.Data["AllSubmitinfo"] = submitinfo

}


func (c *SubmitterController) SubmitProcessing() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))
	//验证token


	//fmt.Println("debug now is ...",claim["username"]) 取值
	c.Data["username"] = "wanlb"
	c.Data["thisName"]="大前端"


	c.getProcessingSubmitData()

	c.TplName = "role/submitter/accessSubmitter.html"
	//c.Layout = "common/layout_home.html"
	c.Layout = "common/layout_home.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["treeviewhtml"] = "role/ops/treeview.html"
}
func (c *SubmitterController) getProcessingSubmitData(){
	var submitinfo []models.Submit_status

	Msql := orm.NewOrm()
	Msql.Using("default")
	//id倒叙排序，前端展示最新提交信息
	Msql.QueryTable("Submit_status").Filter("Submitter",c.Data["username"]).Filter("ApprovalStatus",0).All(&submitinfo,"Id","Dbname","Sqlstatment","Submittime","Approver","ApprovalStatus","Operator","OperatorStatus")
	c.Data["AllSubmitinfo"] = submitinfo

}


