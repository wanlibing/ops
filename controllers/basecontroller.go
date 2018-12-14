package controllers

import (

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"strings"
	"msql/models"
	"k8s.io/client-go/kubernetes"
	"flag"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type BaseController struct {
	beego.Controller
	controllerName string             //当前控制名称
	actionName     string             //当前action名称
	curUser models.User   //记录当前用户信息

}



type ControllerError struct {
	Status   int    `json:"status"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	DevInfo  string `json:"dev_info"`
	MoreInfo string `json:"more_info"`
}

var (
	Err404          = &ControllerError{404, 404, "page not found", "page not found", ""}
	ErrInputData    = &ControllerError{400, 10001, "数据输入错误", "客户端参数错误", ""}
	ErrDatabase     = &ControllerError{500, 10002, "服务器错误", "数据库操作错误", ""}
	ErrDupUser      = &ControllerError{400, 10003, "用户信息已存在", "数据库记录重复", ""}
	ErrNoUser       = &ControllerError{400, 10004, "用户信息不存在", "数据库记录不存在", ""}
	ErrPass         = &ControllerError{400, 10005, "用户信息不存在或密码不正确", "密码不正确", ""}
	ErrNoUserPass   = &ControllerError{400, 10006, "用户信息不存在或密码不正确", "数据库记录不存在或密码不正确", ""}
	ErrNoUserChange = &ControllerError{400, 10007, "用户信息不存在或数据未改变", "数据库记录不存在或数据未改变", ""}
	ErrInvalidUser  = &ControllerError{400, 10008, "用户信息不正确", "Session信息不正确", ""}
	ErrOpenFile     = &ControllerError{500, 10009, "服务器错误", "打开文件出错", ""}
	ErrWriteFile    = &ControllerError{500, 10010, "服务器错误", "写文件出错", ""}
	ErrSystem       = &ControllerError{500, 10011, "服务器错误", "操作系统错误", ""}
	ErrExpired      = &ControllerError{400, 10012, "登录已过期", "验证token过期", ""}
	ErrPermission   = &ControllerError{400, 10013, "没有权限", "没有操作权限", ""}
	Actionsuccess   = &ControllerError{200, 90000, "操作成功", "操作成功", ""}
)


/*
func (c *BaseController) ParseHeaderToken() (t *jwt.Token ) {
	authString := c.Ctx.Input.Header("Authorization")
	beego.Debug("AuthString:", authString)

	fmt.Println("debug base 01")

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
		return nil
	}
	tokenString := kv[1]

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("tokensecret")), nil
	})
	if err != nil {
		beego.Error("Parse token:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// That‘s not even a token
				return nil
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return nil
			} else {
				// Couldn‘t handle this token
				return nil
			}
		} else {
			// Couldn‘t handle this token
			return nil
		}
	}
	if !token.Valid {
		beego.Error("Token invalid:", tokenString)
		return nil
	}
	beego.Debug("Token:", token)

	return token
}
*/
//token来自get请求参数
func (c *BaseController) ParseParamsToken() (t *jwt.Token ) {
	authString := c.GetString("token")
	beego.Debug("AuthString:", authString)

	fmt.Println("debug base 01")

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
		return nil
	}
	tokenString := kv[1]

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("tokensecret")), nil
	})
	if err != nil {
		beego.Error("Parse token:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// That‘s not even a token
				return nil
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return nil
			} else {
				// Couldn‘t handle this token
				return nil
			}
		} else {
			// Couldn‘t handle this token
			return nil
		}
	}
	if !token.Valid {
		beego.Error("Token invalid:", tokenString)
		return nil
	}
	beego.Debug("Token:", token)

	return token
}
//token来自header
// ParseToken parse JWT token in http header.
func (c *BaseController) ParseHeaderToken() (t *jwt.Token, e *ControllerError) {
	authString := c.Ctx.Input.Header("Authorization")
	beego.Debug("AuthString:", authString)

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
		return nil, ErrInputData
	}
	tokenString := kv[1]

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("tokensecret")), nil
	})
	if err != nil {
		beego.Error("Parse token:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// That's not even a token
				return nil, ErrInputData
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return nil, ErrExpired
			} else {
				// Couldn't handle this token
				return nil, ErrInputData
			}
		} else {
			// Couldn't handle this token
			return nil, ErrInputData
		}
	}
	if !token.Valid {
		beego.Error("Token invalid:", tokenString)
		return nil, ErrInputData
	}
	beego.Debug("Token:", token)

	return token, nil
}


//kubernetes 配置
var clientset *kubernetes.Clientset

type k8sNode struct {
	Name string
	Ip string
	Status string
	Cpu int64
	Memory int64
	JoinDate string
	System string
	DockerVersion string
	Pods []string
}

func init()  {
	k8sconfig := flag.String("k8sconfig","./k8sconfig","kubernetes config file path")
	flag.Parse()
	config , err := clientcmd.BuildConfigFromFlags("",*k8sconfig)
	if err != nil {
		log.Println(err)
	}
	clientset , err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("connect k8s success")
	}

}

//获取k8s node 节点



//采取这样的方式，只有从login页面进来才能看到，从其他的页面链接过来看不到，怎办？
//只能走ajax？
