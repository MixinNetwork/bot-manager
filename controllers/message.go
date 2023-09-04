package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/bot-manager/externals"
	"github.com/MixinNetwork/bot-manager/models"
	"github.com/MixinNetwork/bot-manager/session"
	"github.com/astaxie/beego"
)

type MessageController struct {
	beego.Controller
}

type AddMessageReplayType struct {
	Keys     []string `json:"keys"`
	Category string   `json:"category"`
	Data     string   `json:"data"`
	ClientId string   `json:"client_id"`
	ReplayId string   `json:"replay_id"`
}

func (c *MessageController) AddMessageReplay() {
	userId := c.Ctx.Input.GetData("UserId")
	body := c.Ctx.Input.RequestBody
	var messageReplay AddMessageReplayType
	err := json.Unmarshal(body, &messageReplay)
	if err != nil {
		log.Println("AddMessageReplay!!", err)
	}
	if !checkBotManager(userId.(string), messageReplay.ClientId, c.Ctx) {
		return
	}
	if messageReplay.ReplayId == "" {
		messageReplay.ReplayId = bot.UuidNewV4().String()
	}
	for _, key := range messageReplay.Keys {
		models.AddOrUpdateAutoReplayMessage(messageReplay.ReplayId, key, messageReplay.ClientId, messageReplay.Category, messageReplay.Data)
	}

	c.Data["json"] = Resp{Data: "ok"}
	c.ServeJSON()
}

func (c *MessageController) GetMessageReplay() {
	clientId := c.GetString("client_id")
	userId := c.Ctx.Input.GetData("UserId")
	if !checkBotManager(userId.(string), clientId, c.Ctx) {
		return
	}
	data := models.GetAutoReplayMessage(clientId)
	c.Data["json"] = Resp{Data: data}
	c.ServeJSON()
}

func (c *MessageController) DeleteMessageReplay() {
	replayId := c.GetString("replay_id")
	clientId := c.GetString("client_id")
	userId := c.Ctx.Input.GetData("UserId")
	if !checkBotManager(userId.(string), clientId, c.Ctx) {
		return
	}
	models.DeleteAutoReplayMessage(replayId)
	c.Data["json"] = Resp{Data: "ok"}
	c.ServeJSON()
}

func (c *MessageController) GetBroadcast() {
	clientId := c.GetString("client_id")
	userId := c.Ctx.Input.GetData("UserId")
	if !checkBotManager(userId.(string), clientId, c.Ctx) {
		return
	}
	if clientBot := models.CheckUserHasBot(userId.(string), clientId); clientBot == nil {
		err := session.ForbiddenError()
		session.HandleError(c.Ctx, err)
		return
	}
	broadcastList := models.GetBroadcast(clientId)
	c.Data["json"] = Resp{Data: broadcastList}
	c.ServeJSON()
}

func (c *MessageController) PostBroadcast() {
	userId := c.Ctx.Input.GetData("UserId")
	reqModel := new(struct {
		ClientId string `json:"client_id"`
		Category string `json:"category"`
		Data     string `json:"data"`
	})
	err := json.Unmarshal(c.Ctx.Input.RequestBody, reqModel)
	if err != nil {
		log.Println("PostBroadcast Unmarshal!!", err)
	}
	if !checkBotManager(userId.(string), reqModel.ClientId, c.Ctx) {
		return
	}
	var clientBot *models.Bot
	if clientBot = models.CheckUserHasBot(userId.(string), reqModel.ClientId); clientBot == nil {
		err := session.ForbiddenError()
		session.HandleError(c.Ctx, err)
		return
	}
	// 1. 先返回结果
	c.Data["json"] = Resp{Data: "ok"}
	c.ServeJSON()
	// 2. 构建一个原始消息
	originMessageId := bot.UuidNewV4().String()
	// 3. 构建广播消息 并发送
	sendMessages, err := sendBroadcast(clientBot, reqModel.Category, reqModel.Data, false)
	if err != nil {
		log.Println("PostBroadcast sendBroadcast", err)
		return
	}
	err = externals.SendText(clientBot, userId.(string), "发送完毕")
	if err != nil {
		log.Println("PostBroadcast sendBroadcast, 2", err)
		return
	}
	models.AddBroadcast(clientBot.ClientId, userId.(string), originMessageId, reqModel.Category, reqModel.Data)
	// 4. 保存所有消息和原始消息的联系
	for _, message := range sendMessages {
		models.AddBroadcastTmpMessage(reqModel.ClientId, message.MessageId, originMessageId, message.RecipientId, message.ConversationId)
	}
}
func (c *MessageController) DeleteBroadcast() {
	userId := c.Ctx.Input.GetData("UserId")
	clientId := c.GetString("client_id")
	messageId := c.GetString("message_id")
	if !checkBotManager(userId.(string), clientId, c.Ctx) {
		return
	}
	var clientBot *models.Bot
	if clientBot = models.CheckUserHasBot(userId.(string), clientId); clientBot == nil {
		err := session.ForbiddenError()
		session.HandleError(c.Ctx, err)
		return
	}
	messageList := models.GetBroadcastTmpMessage(clientId, messageId)
	var sendMessages []*bot.MessageRequest
	for _, tmp := range messageList {
		str := fmt.Sprintf(`{"message_id":"%s"}`, tmp.MessageId)
		base64Data := base64.StdEncoding.EncodeToString([]byte(str))
		sendMessages = append(sendMessages, &bot.MessageRequest{
			ConversationId: tmp.ConversationId,
			RecipientId:    tmp.UserId,
			MessageId:      bot.UuidNewV4().String(),
			Category:       "MESSAGE_RECALL",
			Data:           base64Data,
		})
	}
	err := externals.SendBatchMessage(sendMessages, clientBot.ClientId, clientBot.SessionId, clientBot.PrivateKey)
	if err != nil {
		return
	}
	models.DeleteBroadcastTmp(clientId, messageId)
	models.DeleteBroadcast(clientId, messageId)
	c.Data["json"] = Resp{Data: "ok"}
	c.ServeJSON()
}

func (c *MessageController) UploadFile() {
	f, h, _ := c.GetFile("file")
	ext := path.Ext(h.Filename)
	if ok := allowExtMap[ext]; !ok {
		session.HandleBadRequestError(c.Ctx)
		return
	}
	defer f.Close()
	obj, err := externals.UploadFile(f, h.Size)
	if err != nil {
		log.Println("UploadFile Err", err)
	}
	c.Data["json"] = Resp{Data: obj}
	c.ServeJSON()
}

func ConnectAllBot() {
	bots := models.GetAllBot()
	for _, userBot := range bots {
		go connectBot(userBot)
	}
	time.Sleep(time.Second)
	go connectWss()
}

type ImageMessagePayload struct {
	AttachmentId string `json:"attachment_id"`
	Size         string `json:"size"`
	Thumbnail    string `json:"thumbnail"`
	Width        string `json:"width"`
	Height       string `json:"height"`
	MimeType     string `json:"mime_type"`
}

var allowExtMap = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}
