package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"github.com/liuzemei/bot-manager/controllers"
	"github.com/liuzemei/bot-manager/db"
	_ "github.com/liuzemei/bot-manager/middleware"
	"github.com/liuzemei/bot-manager/models"
	_ "github.com/liuzemei/bot-manager/routers"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	db.Connect()
	db.AutoMigrate()
	controllers.ConnectAllBot()

	tk := toolbox.NewTask("saveDailyData", "1 0 0 * * *", func() error {
		models.SaveTodayData(-1)
		return nil
	})
	toolbox.AddTask("myTask", tk)
	toolbox.StartTask()
	go func() {
		runtime.SetBlockProfileRate(1) // 开启对阻塞操作的跟踪
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	beego.Run()
}
