package controllers

import (
	"github.com/astaxie/beego"
)

//默认首页登录入口控制器
type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "login.html"
}

