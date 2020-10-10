package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/astaxie/beego"
	"github.com/liuzemei/bot-manager/durable"
	"github.com/liuzemei/bot-manager/externals"
	"github.com/liuzemei/bot-manager/models"
	"github.com/liuzemei/bot-manager/session"
	"github.com/liuzemei/bot-manager/utils"
	"log"
	"path"
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

	// 1. 拿到所有的用户
	userList := models.GetBotUserListById(reqModel.ClientId)
	// 2. 构建一个 原始消息
	originMessageId := bot.UuidNewV4().String()
	models.AddBroadcast(reqModel.ClientId, userId.(string), originMessageId, reqModel.Category, reqModel.Data)
	// 3. 构建所有的消息
	var sendMessages []*bot.MessageRequest
	base64Data := ""
	if reqModel.Category == "PLAIN_TEXT" {
		base64Data = base64.StdEncoding.EncodeToString([]byte(reqModel.Data))
	} else if reqModel.Category == "PLAIN_IMAGE" {
		_msgData, err := json.Marshal(reqModel.Data)
		base64Data = base64.StdEncoding.EncodeToString(_msgData)
		if err != nil {
			return
		}
	}
	for _, user := range userList {
		sendMessages = append(sendMessages, &bot.MessageRequest{
			ConversationId: bot.UniqueConversationId(reqModel.ClientId, user),
			RecipientId:    user,
			MessageId:      bot.UuidNewV4().String(),
			Category:       reqModel.Category,
			Data:           base64Data,
		})
	}
	err = bot.PostMessages(durable.Ctx, sendMessages, clientBot.ClientId, clientBot.SessionId, clientBot.PrivateKey)
	if err != nil {
		log.Println(err)
	}
	// 4. 保存所有消息和原始消息的联系
	for _, message := range sendMessages {
		models.AddBroadcastTmpMessage(reqModel.ClientId, message.MessageId, originMessageId, message.RecipientId, message.ConversationId)
	}
	c.Data["json"] = Resp{Data: "ok"}
	c.ServeJSON()
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
		log.Println(str)
		base64Data := base64.StdEncoding.EncodeToString([]byte(str))
		sendMessages = append(sendMessages, &bot.MessageRequest{
			ConversationId: tmp.ConversationId,
			RecipientId:    tmp.UserId,
			MessageId:      bot.UuidNewV4().String(),
			Category:       "MESSAGE_RECALL",
			Data:           base64Data,
		})
	}
	err := bot.PostMessages(durable.Ctx, sendMessages, clientBot.ClientId, clientBot.SessionId, clientBot.PrivateKey)
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
	for _, bot := range bots {
		go connectBot(bot)
	}
}

type ImageMessagePayload struct {
	AttachmentId string `json:"attachment_id"`
	Size         string `json:"size"`
	Thumbnail    string `json:"thumbnail"`
	Width        string `json:"width"`
	Height       string `json:"height"`
	MimeType     string `json:"mime_type"`
}

func readMessage(messageCome <-chan models.MessengerChannel, hash string) {
	for msg := range messageCome {
		_msg := &msg.Message
		if _msg.Source == "ACKNOWLEDGE_MESSAGE_RECEIPT" {
			models.UpdateMessage(_msg.MessageId, _msg.Status)
			for _, chanel := range models.HashManagerMap[hash] {
				if chanel != nil {
					chanel <- models.RespMessage{
						MessageId: _msg.MessageId,
						Status:    _msg.Status,
						Source:    _msg.Source,
					}
				}
			}
			continue
		}

		// 数据统计
		adminIds := models.GetAdminIdsByBotId(msg.ClientID)
		isManager := false
		for _, id := range adminIds {
			if id == _msg.UserId {
				isManager = true
				break
			}
		}
		if !isManager && _msg.Category == "PLAIN_TEXT" {
			decodeBytes, _ := base64.StdEncoding.DecodeString(_msg.Data)
			decodeString := string(decodeBytes)
			replayData, replayCategory := models.GetAutoReplayMessageByKey(msg.ClientID, decodeString)
			if replayData != "" {
				forwardDashboardMessage(&forwardMessagePropsType{
					Category:         _msg.Category,
					CreatedAt:        utils.FormatTime(_msg.CreatedAt),
					MessageId:        _msg.MessageId,
					Source:           _msg.Source,
					UserId:           _msg.UserId,
					QuoteMessageId:   _msg.QuoteMessageId,
					ConversationId:   _msg.ConversationId,
					Data:             _msg.Data,
					RepresentativeId: _msg.RepresentativeId,
					Status:           _msg.Status,
					UpdatedAt:        utils.FormatTime(_msg.UpdatedAt),
				}, msg.ClientID, msg.SessionID, msg.PrivateKey, hash, false, adminIds)
				var sendMessages []*bot.MessageRequest
				base64Data := base64.StdEncoding.EncodeToString([]byte(replayData))
				messageId := bot.UuidNewV4().String()
				sendMessages = append(sendMessages, &bot.MessageRequest{
					MessageId:      messageId,
					Category:       replayCategory,
					Data:           base64Data,
					ConversationId: _msg.ConversationId,
					RecipientId:    _msg.UserId,
				})
				err := bot.PostMessages(durable.Ctx, sendMessages, msg.ClientID, msg.SessionID, msg.PrivateKey)
				if err != nil {
					log.Println("转发自动回复的消息出问题了。", err)
					continue
				}
				forwardDashboardMessage(&forwardMessagePropsType{
					Category:       replayCategory,
					CreatedAt:      utils.FormatTime(_msg.CreatedAt),
					MessageId:      messageId,
					UserId:         _msg.UserId,
					ConversationId: _msg.ConversationId,
					Data:           base64Data,
					UpdatedAt:      utils.FormatTime(_msg.UpdatedAt),
				}, msg.ClientID, msg.SessionID, msg.PrivateKey, hash, false, adminIds)
				continue
			}
		}

		//  2. 转发给已登录的后端管理的管理员
		forwardDashboardMessage(&forwardMessagePropsType{
			Category:         _msg.Category,
			CreatedAt:        utils.FormatTime(_msg.CreatedAt),
			MessageId:        _msg.MessageId,
			Source:           _msg.Source,
			UserId:           _msg.UserId,
			QuoteMessageId:   _msg.QuoteMessageId,
			ConversationId:   _msg.ConversationId,
			Data:             _msg.Data,
			RepresentativeId: _msg.RepresentativeId,
			Status:           _msg.Status,
			UpdatedAt:        utils.FormatTime(_msg.UpdatedAt),
		}, msg.ClientID, msg.SessionID, msg.PrivateKey, hash, isManager, adminIds)

		// 转发给手机端
		if isManager {
			// 管理员的消息
			//	1. 转发给其他管理员
			//  2. 转发给被 Quote 的用户
			if _msg.QuoteMessageId != "" {
				originMessage := models.GetOriginMessageById(msg.ClientID, _msg.QuoteMessageId)
				if originMessage == nil {
					log.Println("没找到这条消息originMessage，", _msg.QuoteMessageId)
					continue
				}
				forwardMessages := models.GetForwardMessagesByOrigin(msg.ClientID, originMessage.OriginMessageId)
				if forwardMessages == nil {
					log.Println("没找到这条消息forwardMessage", originMessage.OriginMessageId)
					continue
				}
				// 转发给其他管理员，包含quote自己
				var sendMessages []*bot.MessageRequest
				for _, userId := range adminIds {
					if userId == _msg.UserId {
						continue
					}
					conversationId := bot.UniqueConversationId(msg.ClientID, userId)
					sendMessages = append(sendMessages, &bot.MessageRequest{
						MessageId:        bot.UuidNewV4().String(),
						Category:         _msg.Category,
						Data:             _msg.Data,
						QuoteMessageId:   forwardMessages[userId].MessageId,
						ConversationId:   conversationId,
						RecipientId:      userId,
						RepresentativeId: _msg.UserId,
					})
				}
				// 转发给被 Quote 的那个用户
				sendMessages = append(sendMessages, &bot.MessageRequest{
					ConversationId: originMessage.ConversationId,
					RecipientId:    originMessage.RecipientId,
					MessageId:      bot.UuidNewV4().String(),
					Category:       _msg.Category,
					Data:           _msg.Data,
					QuoteMessageId: originMessage.OriginMessageId,
				})
				err := bot.PostMessages(durable.Ctx, sendMessages, msg.ClientID, msg.SessionID, msg.PrivateKey)
				if err != nil {
					log.Println("转发管理员消息出问题了。", err)
					continue
				}
				for _, message := range sendMessages {
					models.AddForwardMessage(models.ForwardMessage{
						ClientId:        msg.ClientID,
						MessageId:       message.MessageId,
						AdminId:         _msg.UserId,
						RecipientId:     message.RecipientId,
						OriginMessageId: originMessage.MessageId,
						ConversationId:  message.ConversationId,
						AdminMessageId:  _msg.MessageId,
						CreatedAt:       utils.FormatTime(_msg.CreatedAt),
					})
				}
			}
		} else {
			// 其他人员的消息
			//	1. 直接转发给管理员的 Messenger
			var sendMessages []*bot.MessageRequest
			for _, userId := range adminIds {
				conversationId := bot.UniqueConversationId(msg.ClientID, userId)
				messageId := bot.UuidNewV4().String()
				sendMessages = append(sendMessages, &bot.MessageRequest{
					ConversationId:   conversationId,
					RecipientId:      userId,
					RepresentativeId: _msg.UserId,
					MessageId:        messageId,
					Category:         _msg.Category,
					Data:             _msg.Data,
				})
				models.AddForwardMessage(models.ForwardMessage{
					ClientId:        msg.ClientID,
					MessageId:       messageId,
					OriginMessageId: _msg.MessageId,
					RecipientId:     _msg.UserId,
					ConversationId:  _msg.ConversationId,
					AdminId:         userId,
					CreatedAt:       utils.FormatTime(_msg.CreatedAt),
				})
			}
			err := bot.PostMessages(durable.Ctx, sendMessages, msg.ClientID, msg.SessionID, msg.PrivateKey)
			if err != nil {
				log.Println("转发普通消息出问题了。", err)
			}
		}
	}
}

var allowExtMap = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}
