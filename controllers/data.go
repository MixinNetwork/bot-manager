package controllers

import (
	"github.com/astaxie/beego"
	"github.com/liuzemei/bot-manager/models"
)

type DataController struct {
	beego.Controller
}

func (c *DataController) Get() {
	clientId := c.Ctx.Input.Query("client_id")
	resp := models.GetDailyData(clientId)
	c.Data["json"] = Resp{Data: resp}
	c.ServeJSON()
}
