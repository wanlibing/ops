package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/config"
	"fmt"
)

func newDbConfig() string{
	dbconfig , err := config.NewConfig("ini","conf/db.conf")
	if err != nil{
		return err.Error()
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		dbconfig.String("mysql::db_user"),
		dbconfig.String("mysql::db_pwd"),
		dbconfig.String("mysql::db_host"),
		dbconfig.String("mysql::db_port"),
		dbconfig.String("mysql::db_name"))
}

//var Msql  = orm.NewOrm()

func init()  {
	connstr := newDbConfig()
	orm.RegisterModel(new(User),new(Submit_status),new(SqlPushRelationship)) //顺序無所謂
	orm.RegisterDriver("mysql",orm.DRMySQL)
	//orm.RegisterDataBase("default","mysql","wanlb:123456a@tcp(192.168.138.128:3306)/msql?charset=utf8")
	orm.RegisterDataBase("default","mysql",connstr)  //MaxIdleConns 可以设置最大空闲连接数,MaxOpenConns 最大连接数
	//调试模式
	orm.Debug = true
	//自动建表
	orm.RunSyncdb("default",false,true)
}