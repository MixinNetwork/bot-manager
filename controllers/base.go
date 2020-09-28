package controllers

import "github.com/astaxie/beego"

type Resp struct {
	Data interface{} `json:"data"`
}

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Options() {
	c.Data["json"] = map[string]interface{}{"status": 200, "message": "ok"}
	c.ServeJSON()
}
