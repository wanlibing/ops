package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"msql/models"
)

//默认首页登录入口控制器
type CheckerController struct {
	beego.Controller
}

func (c *CheckerController) Get() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))
	c.getCheckerData()
	c.Data["username"] = "qiushu"    //仅作测试
	c.Data["thisName"]="大前端"
	c.TplName = "role/checker/checker.html"
	c.Layout = "common/layout_home.html"
}

//采用排序来控制在页面上的展示最新的任务，并用js实现是否能撤销操作
func (c *CheckerController) getCheckerData(){
	var submitinfo []models.Submit_status
	Msql := orm.NewOrm()
	Msql.Using("default")
	//id倒叙排序，前端展示最新提交信息
	Msql.QueryTable("Submit_status").Filter("Approver","qiushu").All(&submitinfo,"Id","Dbname","Sqlstatment","Submittime","Approver","ApprovalStatus","Operator","OperatorStatus")
	c.Data["AllSubmitinfo"] = submitinfo

}

