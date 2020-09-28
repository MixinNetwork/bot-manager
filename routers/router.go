// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/liuzemei/bot-manager/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("*", &controllers.BaseController{}, "OPTIONS:Options"),
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSRouter("/", &controllers.UserController{}, "post:Post"),
			beego.NSRouter("login", &controllers.UserController{}, "get:Login"),
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/data",
			beego.NSRouter("/", &controllers.DataController{}, "get:Get"),
			beego.NSInclude(
				&controllers.DataController{},
			),
		),
		beego.NSNamespace("/bot",
			beego.NSRouter("/", &controllers.BotController{}, "POST:Add;GET:Get"),
			beego.NSInclude(
				&controllers.BotController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
