package models

//用户表
type User struct {
	Id int
	Name string
	Password string   //服务端如何存储密码?
	Department string  //部门
	Roleid int //1:ops ,2:checker, 3:submitter
}

