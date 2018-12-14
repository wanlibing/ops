package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"msql/models"
	"strconv"
)

//默认首页登录入口控制器
type UserController struct {
	beego.Controller
}

//用户管理首页
func (c *UserController) SqlUserManagement() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))

	c.TplName = "usermanager/sqlmanager.html"
	//c.Layout = "common/layout_home.html"
	c.Layout = "common/layout_home.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["treeviewhtml"] = "role/ops/treeview.html"
}

//在线编辑用户信息
func (c *UserController) SqlUserUpdate() {
	Msql := orm.NewOrm()
	Msql.Using("default")
	username := c.GetString("username")
	department := c.GetString("department")
	roleid := c.GetString("roleid")
	Msql.QueryTable("User").Filter("Name",username).Update(orm.Params{"Department":department,"Roleid":roleid})

	result := struct {
		Status string
	}{"ok"}
	c.Data["json"] = &result
	c.ServeJSON()


}

//新增用户
func (c *UserController) SqlUserAdd() {
	Msql := orm.NewOrm()
	Msql.Using("default")
	username := c.GetString("username")
	password := c.GetString("password")
	department := c.GetString("department")
	roleid,_ := strconv.Atoi(c.GetString("roleid"))

	inserUser := models.User{Name:username,Password:password,Department:department,Roleid:roleid}
	Msql.Insert(&inserUser)

	result := struct {
		Status string
	}{"ok"}
	c.Data["json"] = &result
	c.ServeJSON()


}
