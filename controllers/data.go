package controllers

import (
	"github.com/astaxie/beego"
	"github.com/liuzemei/bot-manager/models"
)

type DataController struct {
	beego.Controller
}

func (c *DataController) Get() {
	userId := c.Ctx.Input.GetData("UserId")
	resp := models.GetDailyData(userId.(string))
	c.Data["json"] = Resp{Data: resp}
	c.ServeJSON()
}
