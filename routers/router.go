// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/MixinNetwork/bot-manager/controllers"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("*", &controllers.BaseController{}, "OPTIONS:Options"),
		beego.NSNamespace("/user",
			beego.NSRouter("/", &controllers.UserController{}, "GET:Get;POST:Post;PUT:Put"),
			beego.NSRouter("/login", &controllers.UserController{}, "GET:Login"),
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/data",
			beego.NSRouter("/", &controllers.DataController{}, "GET:Get"),
			beego.NSInclude(
				&controllers.DataController{},
			),
		),
		beego.NSNamespace("/bot",
			beego.NSRouter("/", &controllers.BotController{}, "GET:Get;POST:Add"),
			beego.NSRouter("/favorite", &controllers.BotController{}, "GET:FavoriteGet;POST:FavoriteAdd;DELETE:FavoriteDel"),
			beego.NSInclude(
				&controllers.BotController{},
			),
		),
		beego.NSNamespace("/message",
			beego.NSRouter("/uploadFile", &controllers.MessageController{}, "POST:UploadFile"),
			beego.NSRouter("/replay", &controllers.MessageController{}, "GET:GetMessageReplay;POST:AddMessageReplay;DELETE:DeleteMessageReplay"),
			beego.NSRouter("/broadcast", &controllers.MessageController{}, "GET:GetBroadcast;POST:PostBroadcast;DELETE:DeleteBroadcast"),
			beego.NSInclude(
				&controllers.MessageController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
