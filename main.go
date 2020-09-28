package main

import (
	"github.com/astaxie/beego"
	"github.com/liuzemei/bot-manager/controllers"
	"github.com/liuzemei/bot-manager/db"
	_ "github.com/liuzemei/bot-manager/middleware"
	_ "github.com/liuzemei/bot-manager/routers"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	db.Connect()
	db.AutoMigrate()
	controllers.ConnectAllBot()

	beego.Run()
}
