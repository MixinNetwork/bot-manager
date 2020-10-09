package controllers

import (
	"encoding/json"
	botApi "github.com/MixinNetwork/bot-api-go-client"
	"github.com/astaxie/beego"
	"github.com/liuzemei/bot-manager/durable"
	"github.com/liuzemei/bot-manager/externals"
	"github.com/liuzemei/bot-manager/models"
	"github.com/liuzemei/bot-manager/session"
	"log"
)

type BotController struct {
	beego.Controller
}

type addBotReq struct {
	ClientId   string `json:"client_id"`
	SessionId  string `json:"session_id"`
	PrivateKey string `json:"private_key"`
}

func (c *BotController) Add() {
	if err := recover(); err != nil {
		session.HandleBadRequestError(c.Ctx)
	}
	userId := c.Ctx.Input.GetData("UserId").(string)
	body := c.Ctx.Input.RequestBody

	var bot addBotReq
	if err := json.Unmarshal(body, &bot); err != nil {
		session.HandleBadRequestError(c.Ctx)
		return
	}
	bot.PrivateKey = externals.HandlePrivateKey(bot.PrivateKey)
	token, err := botApi.SignAuthenticationToken(bot.ClientId, bot.SessionId, bot.PrivateKey, "GET", "/me", "")
	if err != nil {
		log.Println("/controllers/bot.Add SignAuthenticationToken error", err)
		session.HandleBadRequestError(c.Ctx)
		return
	}
	user, err := botApi.UserMe(durable.Ctx, token)
	if err != nil {
		log.Println("/controllers/bot.Add UserMe error", err)
		session.HandleBadRequestError(c.Ctx)
		return
	}
	if user.UserId != "" {
		models.AddOrUpdateUserBotItem(userId, bot.ClientId, bot.SessionId, bot.PrivateKey)
		models.AddOrUpdateBotItem(bot.ClientId, bot.SessionId, bot.PrivateKey, user.FullName, user.IdentityNumber, user.AvatarURL)
	}

	go connectBot(models.UserBot{
		ClientId:   bot.ClientId,
		SessionId:  bot.SessionId,
		PrivateKey: bot.PrivateKey,
		Hash:       models.Sha256Hash(bot.ClientId, bot.SessionId, bot.PrivateKey),
	})

	c.Data["json"] = Resp{Data: models.BotInfoRes{
		ClientId:       bot.ClientId,
		FullName:       user.FullName,
		IdentityNumber: user.IdentityNumber,
		AvatarURL:      user.AvatarURL,
	}}
	c.ServeJSON()
}

func (c *BotController) Get() {
	userId := c.Ctx.Input.GetData("UserId")
	bots := models.GetUserBotByUserId(userId.(string))
	c.Data["json"] = Resp{Data: bots}
	c.ServeJSON()
}
