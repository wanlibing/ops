package controllers

//默认首页登录入口控制器
type OpsenvController struct {
	BaseController
}

func (c *OpsenvController) Get() {
	c.TplName = "opsenv/baseEnv.html"
	//c.Layout = "common/layout_home.html"
	c.Layout = "common/layout_home.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["treeviewhtml"] = "role/ops/treeview.html"
}

