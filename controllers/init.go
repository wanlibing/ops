package controllers

//定义chan,用户保存推送信息
//是不是那个数据库都可以删除？
//必须做超时机制

type urgedmsg struct {
	Pusher string
	Sqlid int
	Receiver string
}

var urgedmeschan chan urgedmsg

func init()  {
	urgedmeschan = make(chan urgedmsg,100)
}
