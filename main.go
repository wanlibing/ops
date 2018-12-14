package main

import (
	_ "msql/routers"
	_ "msql/models"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

