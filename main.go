package main

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/MixinNetwork/bot-manager/controllers"
	"github.com/MixinNetwork/bot-manager/db"
	_ "github.com/MixinNetwork/bot-manager/middleware"
	"github.com/MixinNetwork/bot-manager/models"
	_ "github.com/MixinNetwork/bot-manager/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
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
