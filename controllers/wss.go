package controllers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/bot-manager/durable"
	"github.com/MixinNetwork/bot-manager/externals"
	"github.com/MixinNetwork/bot-manager/middleware"
	"github.com/MixinNetwork/bot-manager/models"
	"github.com/MixinNetwork/bot-manager/session"
	"github.com/MixinNetwork/bot-manager/utils"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func connectBot(botInfo models.UserBot) {
	if models.HashMessengerMap[botInfo.Hash] != nil {
		return
	}
	var messageTransfer = make(chan models.MessengerChannel)
	models.HashMessengerMap[botInfo.Hash] = messageTransfer
	models.HashManagerMap[botInfo.Hash] = make(map[string]chan models.RespMessage)
	go readMessage(messageTransfer, botInfo.Hash)
	err := externals.StartWebSockets(botInfo.ClientId, botInfo.SessionId, botInfo.PrivateKey, botInfo.Hash)
	if err != nil {
		models.DeleteBotItem(botInfo.ClientId)
		if models.HashMessengerMap[botInfo.Hash] != nil {
			close(models.HashMessengerMap[botInfo.Hash])
			delete(models.HashMessengerMap, botInfo.Hash)
		}
	}
}

func init() {
	go func() {
		http.ListenAndServe(":9099", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if !strings.HasPrefix(r.URL.RawQuery, "token=") {
				errMessage := session.AuthorizationError()
				resp, err := json.Marshal(errMessage)
				if err != nil {
					log.Println("ListenAndServe", err)
				}
				w.Write(resp)
				return
			}
			userId, err := middleware.Parse(r.URL.RawQuery[6:])
			if err != nil || userId == "" {
				// handle error
				return
			}
			hashes := models.GetUserBotHashByUserId(userId)
			if len(hashes) == 0 {
				return
			}
			for _, hash := range hashes {
				var messageTransfer = make(chan models.RespMessage)
				models.HashManagerMap[hash][userId] = messageTransfer
				go wssReadMessage(messageTransfer, conn)
			}
			if models.ChangeBotWss[userId] == nil {
				changeBotChannel := make(chan string, 1)
				models.ChangeBotWss[userId] = changeBotChannel
			}
			go changeWssConnect(models.ChangeBotWss[userId], conn, userId)
			go func() {
				defer conn.Close()
				defer func() {
					if r := recover(); r != nil {
						log.Println("/controllers/wss go func error", r)
					}
				}()
				for {
					msg, op, err := wsutil.ReadClientData(conn)
					if op != 1 || err != nil {
						for hash, hashMap := range models.HashManagerMap {
							for _userId, channel := range hashMap {
								if _userId == userId && channel != nil {
									close(models.HashManagerMap[hash][userId])
									delete(models.HashManagerMap[hash], userId)
								}
							}
						}
						if models.ChangeBotWss[userId] != nil {
							close(models.ChangeBotWss[userId])
							delete(models.ChangeBotWss, userId)
						}
						return
					}
					if string(msg) == "ping" {
						if err := writeMessage(conn, []byte("pong")); err != nil {
							log.Println("/controllers/wss writeMessage", err)
						}
					} else {
						err := handleUserMessage(conn, msg, userId)
						if err != nil {
							log.Println("处理消息错误", err)
						}
					}
				}
			}()
		}))
	}()
}

func handleUserMessage(conn io.Writer, msg []byte, userId string) error {
	var message models.RespMessage
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Println("/controllers/wss handleUserMessage", err)
	}
	if message.CreatedAt != "" {
		var messages []models.RespMessage
		// 发送了 初始化的消息
		if message.ClientId != "" {
			messagesClient := models.GetAllMessagesByBotId(message.ClientId)
			if messagesClient != nil {
				messages = messagesClient
			}
		}
		messagesDate := models.GetAllMessagesByUserId(userId, message.CreatedAt)
		if messagesDate != nil {
			if messages == nil {
				messages = messagesDate
			} else {
				messageIdMap := map[string]bool{}
				for _, respMessage := range messages {
					messageIdMap[respMessage.MessageId] = true
				}
				for _, respMessage := range messagesDate {
					if !messageIdMap[respMessage.MessageId] {
						messages = append(messages, respMessage)
					}
				}
			}
		}

		if messages != nil {
			byteMsg, _ := json.Marshal(messages)
			if err := writeMessage(conn, byteMsg); err != nil {
				return err
			}
		}
	} else {
		// 发送了正常消息
		var msg []byte
		var messagePre *bot.MessageRequest
		if message.Category == "PLAIN_TEXT" {
			base64Data := base64.StdEncoding.EncodeToString([]byte(message.Data))
			messagePre = &bot.MessageRequest{
				ConversationId: bot.UniqueConversationId(message.ClientId, message.RecipientId),
				RecipientId:    message.RecipientId,
				MessageId:      bot.UuidNewV4().String(),
				Category:       message.Category,
				Data:           base64Data,
			}
		} else if message.Category == "PLAIN_IMAGE" {
			_msgData, err := json.Marshal(message.Data)
			base64Data := base64.StdEncoding.EncodeToString(_msgData)
			if err != nil {
				return err
			}
			messagePre = &bot.MessageRequest{
				ConversationId: bot.UniqueConversationId(message.ClientId, message.RecipientId),
				RecipientId:    message.RecipientId,
				MessageId:      bot.UuidNewV4().String(),
				Category:       message.Category,
				Data:           base64Data,
			}
		} else {
			log.Println("不支持的消息格式")
			return nil
		}
		msg, err = json.Marshal(messagePre)
		if err != nil {
			return err
		}
		botInfo := models.GetBotById(message.ClientId)
		accessToken, err := bot.SignAuthenticationToken(botInfo.ClientId, botInfo.SessionId, botInfo.PrivateKey, "POST", "/messages", string(msg))
		if err != nil {
			return err
		}
		body, err := bot.Request(durable.Ctx, "POST", "/messages", msg, accessToken)
		if err != nil {
			return err
		}
		var resp struct {
			Data  bot.MessageView `json:"data"`
			Error bot.Error       `json:"error"`
		}
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return err
		}
		if resp.Error.Code > 0 {
			return resp.Error
		}
		// 收到管理员发来的消息
		// 1. 转发到manager (自己发送的消息)
		// 2. 转发给其他管理员的 messenger。
		data := resp.Data
		adminIds := models.GetAdminIdsByBotId(botInfo.ClientId)
		forwardDashboardMessage(&models.Message{
			Category:         data.Category,
			CreatedAt:        utils.FormatTime(data.CreatedAt),
			MessageId:        data.MessageId,
			Source:           data.Source,
			UserId:           message.RecipientId,
			QuoteMessageId:   data.QuoteMessageId,
			AdminId:          message.UserId,
			ConversationId:   data.ConversationId,
			Data:             data.Data,
			RepresentativeId: data.RepresentativeId,
			Status:           "pending",
			UpdatedAt:        utils.FormatTime(data.UpdatedAt),
		}, botInfo.ClientId, botInfo.SessionId, botInfo.PrivateKey, botInfo.Hash, true, adminIds)

	}
	return nil
}

func forwardDashboardMessage(msg *models.Message, clientId, sessionId, privateKey, hash string, isAdmin bool, adminIds []string) {
	if msg.UserId == "" {
		return
	}
	decodeBytes, _ := base64.StdEncoding.DecodeString(msg.Data)
	msg.Status = strings.ToLower(msg.Status)
	status := msg.Status
	userId := msg.UserId
	data := msg.Data
	switch msg.Category {
	case "PLAIN_TEXT":
		data = string(decodeBytes)
	case "PLAIN_IMAGE":
		var imageMessagePayload ImageMessagePayload
		_ = json.Unmarshal(decodeBytes, &imageMessagePayload)
		att, err := bot.AttachmentShow(durable.Ctx, clientId, sessionId, privateKey, imageMessagePayload.AttachmentId)
		if err != nil {
			return
		}
		data = att.ViewURL
	}
	if msg.AdminId != "" {
		// 来自 网页的消息
		// 1. 发送给其他管理员
		originMessage := models.GetLastMessageByRecipientId(clientId, msg.UserId)
		if originMessage == nil || originMessage.OriginMessageId == "" {
			log.Println("/controllers/wss forwardDashboardMessage msg.AdminId!='' 获取originMessage失败", msg.UserId)
			return
		}
		messages := models.GetForwardMessagesByOrigin(clientId, originMessage.OriginMessageId)
		var sendMessages []*bot.MessageRequest
		for _, id := range adminIds {
			if id == msg.AdminId {
				continue
			}
			conversationId := bot.UniqueConversationId(clientId, id)
			sendMessages = append(sendMessages, &bot.MessageRequest{
				ConversationId:   conversationId,
				RecipientId:      id,
				MessageId:        bot.UuidNewV4().String(),
				Category:         msg.Category,
				Data:             msg.Data,
				RepresentativeId: msg.AdminId,
				QuoteMessageId:   messages[id].MessageId,
			})
		}
		status = "pending"
		if len(sendMessages) > 0 {
			err := externals.SendBatchMessage(sendMessages, clientId, sessionId, privateKey)
			if err != nil {
				log.Println("/controllers/wss forwardDashboardMessage, 发送给其他管理员", err)
			}
		}
	} else if isAdmin {
		// 来自 客户端的 管理员 消息
		// 1.
		originMessage := models.GetOriginMessageById(clientId, msg.QuoteMessageId)
		if originMessage == nil || originMessage.RecipientId == "" {
			log.Println("/controllers/wss forwardDashboardMessage 获取originMessage失败", msg.MessageId)
			return
		} else {
			GetMessageUserAutoUpdate(userId, clientId)
			userId = originMessage.RecipientId
		}
		status = "read"
	}

	userInfo, err := GetMessageUserAutoUpdate(userId, clientId)
	if err != nil || userInfo == nil {
		log.Println("/controllers/wss forwardDashboardMessage GetMessageUserAutoUpdate", err)
	}
	// 发送给后台的管理员
	for _, chanel := range models.HashManagerMap[hash] {
		if chanel != nil {
			var _m models.RespMessage
			_m.Message = *msg
			_m.ClientId = clientId
			_m.IdentityNumber = userInfo.IdentityNumber
			_m.UserId = userId
			_m.FullName = userInfo.FullName
			_m.AvatarURL = userInfo.AvatarURL
			_m.Data = data
			_m.Status = status

			chanel <- _m
		}
	}
	models.UpdateClientMessageById(clientId, userInfo, msg, status)
}

type ackMessageType struct {
	MessageId string `json:"message_id"`
	Status    string `json:"status"`
	Source    string `json:"source"`
}

func ackMessage(conn io.Writer, messageId, status, source string) {
	status = strings.ToLower(status)
	models.UpdateMessage(messageId, status)
	byteMsg, _ := json.Marshal(ackMessageType{
		MessageId: messageId,
		Status:    status,
		Source:    source,
	})
	writeMessage(conn, byteMsg)
}

func wssReadMessage(managerMessageCome chan models.RespMessage, conn io.Writer) {
	log.Println(managerMessageCome)
	for msg := range managerMessageCome {
		if msg.Source == "ACKNOWLEDGE_MESSAGE_RECEIPT" {
			ackMessage(conn, msg.MessageId, msg.Status, msg.Source)
			continue
		}
		byteMsg, err := json.Marshal(msg)
		if err != nil {
			log.Println("消息转发解析出问题了！", err)
		}
		writeMessage(conn, byteMsg)
	}
}
func changeWssConnect(hashChannel chan string, conn io.Writer, userId string) {
	for hash := range hashChannel {
		var messageTransfer = make(chan models.RespMessage)
		models.HashManagerMap[hash][userId] = messageTransfer
		go wssReadMessage(messageTransfer, conn)
	}
}

func writeMessage(conn io.Writer, message []byte) error {
	if err := wsutil.WriteServerMessage(conn, 1, message); err != nil {
		return err
	}
	return nil
}
