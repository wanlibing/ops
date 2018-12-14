package routers

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

var FilterToken = func(ctx *context.Context) {

	authString := ctx.Input.Header("Authorization")
	if  authString ==""  {
		authString = ctx.GetCookie("token")
	}

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
	}
	tokenString := kv[1]

	// Parse token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("tokensecret")), nil
	})
	if !token.Valid {
		beego.Error("Token invalid:", tokenString)
		ctx.Redirect(302,"/login")
	}

}

//过滤器使用
func init()  {
	beego.InsertFilter("/home/submitter/processing",beego.BeforeRouter,FilterToken)
}


