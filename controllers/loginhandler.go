package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/astaxie/beego/orm"
	"msql/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"time"
)

type LoginController struct {
	BaseController
}

type Token struct {
	Token string
}

//验证用户登录
func (c *LoginController) LoginStatus() {
	username ,password := c.GetString("username"),c.GetString("password")
	//获取服务端密码
	var uerinfo models.User
	var status string
	Msql := orm.NewOrm()
	Msql.Using("default")
	err := Msql.QueryTable("user").Filter("name",strings.TrimSpace(username)).One(&uerinfo,"password","roleid","name")
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		fmt.Printf("Returned Multi Rows Not One")
	}
	if err == orm.ErrNoRows {
		// 没有找到记录
		fmt.Printf("Not row found")

	}
	if strings.TrimSpace(password) == strings.TrimSpace(uerinfo.Password)  {
		//登录成功
		//设置token到cookie
		token := jwt.New(jwt.SigningMethodHS256)
		claims := make(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()  //token过期时间：1小时
		claims["iat"] = time.Now().Unix()
		claims["username"] = uerinfo.Name   //将用户名和角色ID存入到TOKEN
		claims["roleid"] = uerinfo.Roleid
		token.Claims = claims
		//fmt.Println("app token config secret is ",beego.AppConfig.String("tokensecret"))
		//签证
		tokenString, _ :=  token.SignedString([]byte(beego.AppConfig.String("tokensecret")))

		c.Ctx.SetCookie("token","Bearer " + tokenString)

		//c.curUser = models.User{Name:uerinfo.Name,Roleid:uerinfo.Roleid}  如何使用


		status = "successed"
	} else {
		status = "failed"
	}
	fmt.Println("login status is ",status)
	//返回登录状态，后续要带roleid跳转页面
	result := struct {
		Status string
		RoleId int
	}{status,uerinfo.Roleid}
	c.Data["json"] = &result
	c.ServeJSON()

}
