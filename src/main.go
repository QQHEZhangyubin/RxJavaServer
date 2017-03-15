package main

import (
	_ "./routers"
	"./models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	models.InitTable();
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, false)
	models.AutoCreateData()
	beego.Run()

}

