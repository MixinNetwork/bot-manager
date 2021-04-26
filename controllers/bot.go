package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	botApi "github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/bot-manager/durable"
	"github.com/MixinNetwork/bot-manager/externals"
	"github.com/MixinNetwork/bot-manager/models"
	"github.com/MixinNetwork/bot-manager/session"
	"github.com/astaxie/beego"
)

type BotController struct {
	beego.Controller
}

type favoriteApp struct {
	Type      string `json:"type"`
	UserId    string `json:"user_id"`
	AppId     string `json:"app_id"`
	CreatedAt string `json:"created_at"`
}

func (c *BotController) Add() {
	if err := recover(); err != nil {
		session.HandleBadRequestError(c.Ctx)
	}
	userId := c.Ctx.Input.GetData("UserId").(string)
	body := c.Ctx.Input.RequestBody

	var bot models.UserBot
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

	c.Data["json"] = Resp{Data: models.Bot{
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

func (c *BotController) FavoriteGet() {
	userId := c.Ctx.Input.GetData("UserId")
	clientId := c.GetString("client_id")
	if !checkBotManager(userId.(string), clientId, c.Ctx) {
		return
	}
	var clientBot *models.Bot
	if clientBot = models.CheckUserHasBot(userId.(string), clientId); clientBot == nil {
		err := session.ForbiddenError()
		session.HandleError(c.Ctx, err)
		return
	}
	uri := fmt.Sprintf("/users/%s/apps/favorite", clientId)
	accessToken, err := botApi.SignAuthenticationToken(clientBot.ClientId, clientBot.SessionId, clientBot.PrivateKey, "GET", uri, "")
	if err != nil {
		return
	}
	body, err := botApi.Request(durable.Ctx, "GET", uri, []byte{}, accessToken)
	if err != nil {
		return
	}
	var _resp struct {
		Data  []favoriteApp `json:"data"`
		Error botApi.Error  `json:"error"`
	}
	err = json.Unmarshal(body, &_resp)
	if err != nil {
		return
	}
	var userList []string
	for _, app := range _resp.Data {
		userList = append(userList, app.AppId)
	}
	resp := make([]*models.User, 0)
	if userList != nil {
		resp = models.GetUserByIds(userList)
	}
	c.Data["json"] = Resp{Data: resp}
	c.ServeJSON()
}

func (c *BotController) FavoriteAdd() {
	userId := c.Ctx.Input.GetData("UserId")
	clientId := c.GetString("client_id")
	id := c.GetString("id")
	if !checkBotManager(userId.(string), clientId, c.Ctx) {
		return
	}
	var clientBot *models.Bot
	if clientBot = models.CheckUserHasBot(userId.(string), clientId); clientBot == nil {
		err := session.ForbiddenError()
		session.HandleError(c.Ctx, err)
		return
	}
	user, err := botApi.SearchUser(durable.Ctx, id, clientBot.ClientId, clientBot.SessionId, clientBot.PrivateKey)
	if err != nil || user == nil || user.UserId == "" {
		session.HandleBadRequestError(c.Ctx)
		return
	}
	models.AddUser(&models.User{
		UserId:         user.UserId,
		FullName:       user.FullName,
		IdentityNumber: user.IdentityNumber,
		AvatarURL:      user.AvatarURL,
		AccessToken:    "",
		CreatedAt:      user.CreatedAt,
	})
	uri := fmt.Sprintf("/apps/%s/favorite", user.UserId)
	accessToken, err := botApi.SignAuthenticationToken(clientBot.ClientId, clientBot.SessionId, clientBot.PrivateKey, "POST", uri, "")
	if err != nil {
		return
	}
	body, err := botApi.Request(durable.Ctx, "POST", uri, []byte{}, accessToken)
	if err != nil {
		session.HandleBadRequestError(c.Ctx)
		return
	}
	var resp struct {
		Data  favoriteApp  `json:"data"`
		Error botApi.Error `json:"error"`
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}
	c.Data["json"] = Resp{Data: resp.Data}
	c.ServeJSON()
}

func (c *BotController) FavoriteDel() {
	userId := c.Ctx.Input.GetData("UserId")
	clientId := c.GetString("client_id")
	id := c.GetString("id")
	if !checkBotManager(userId.(string), clientId, c.Ctx) {
		return
	}
	var clientBot *models.Bot
	if clientBot = models.CheckUserHasBot(userId.(string), clientId); clientBot == nil {
		err := session.ForbiddenError()
		session.HandleError(c.Ctx, err)
		return
	}
	uri := fmt.Sprintf("/apps/%s/unfavorite", id)
	accessToken, err := botApi.SignAuthenticationToken(clientBot.ClientId, clientBot.SessionId, clientBot.PrivateKey, "POST", uri, "")
	if err != nil {
		return
	}
	_, err = botApi.Request(durable.Ctx, "POST", uri, []byte{}, accessToken)
	if err != nil {
		session.HandleBadRequestError(c.Ctx)
		return
	}
	c.Data["json"] = Resp{Data: "ok"}
	c.ServeJSON()
}
