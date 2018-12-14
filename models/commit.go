package models

//提交关系表,所有用户都提交到此表
type Submit_status struct {
	Id int   //任务ID sqlid
	Dbname string  //数据库
	Sqlstatment string //语句（后续存储在mongodb）
	Submittime string  //提交时间
	Submitter string //提交人
	Approver string    //审核负责人
	ApprovalStatus int   //审核状态 0：审批中 ； 1：审批通过 2.审批不通过 ；
	Operator string     //执行人
	OperatorStatus int  //执行状态 0：未执行 1：已执行 2.拒绝执行
}


//审核执行关系表
type SqlPushRelationship struct {
	Id int     //beego 必须要求主键
	Pusher string
	SqlId int
	Receiver string
}
