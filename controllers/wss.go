package controllers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/liuzemei/bot-manager/durable"
	"github.com/liuzemei/bot-manager/middleware"
	"github.com/liuzemei/bot-manager/models"
	"github.com/liuzemei/bot-manager/session"
	"github.com/liuzemei/bot-manager/utils"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	go func() {
		http.ListenAndServe(":9099", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if !strings.HasPrefix(r.URL.RawQuery, "token=") {
				errMessage := session.AuthorizationError()
				resp, err := json.Marshal(errMessage)
				if err != nil {
					log.Println(err)
				}
				w.Write(resp)
				return
			}
			userId, err := middleware.Parse(r.URL.RawQuery[6:])
			log.Println(userId)

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
				for {
					msg, op, err := wsutil.ReadClientData(conn)
					if err != nil {
						// handle error
					}
					if op != 1 {
						for hash, hashMap := range models.HashManagerMap {
							for _userId, channel := range hashMap {
								if _userId == userId && channel != nil {
									delete(models.HashManagerMap[hash], userId)
								}
							}
						}
						delete(models.ChangeBotWss, userId)
						return
					}
					if string(msg) == "ping" {
						if err := writeMessage(conn, []byte("pong")); err != nil {
							log.Println(err)
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
	log.Printf("%#v\n", message)
	if err != nil {
		log.Println(err)
	}
	if message.CreatedAt != "" {
		// 发送了 初始化的消息
		messages := models.GetAllMessagesByUserId(userId, message.CreatedAt)
		if messages != nil {
			byteMsg, _ := json.Marshal(messages)
			if err := writeMessage(conn, byteMsg); err != nil {
				return err
			}
		}
	} else {
		// 发送了正常消息
		botInfo := models.GetBotById(message.ClientId)
		base64Data := base64.StdEncoding.EncodeToString([]byte(message.Data))
		messagePre := &bot.MessageRequest{
			ConversationId: bot.UniqueConversationId(message.ClientId, message.RecipientId),
			RecipientId:    message.RecipientId,
			MessageId:      bot.UuidNewV4().String(),
			Category:       message.Category,
			Data:           base64Data,
		}
		msg, err := json.Marshal(messagePre)
		if err != nil {
			return err
		}
		accessToken, err := bot.SignAuthenticationToken(botInfo.ClientId, botInfo.SessionId, botInfo.PrivateKey, "POST", "/messages", string(msg))
		if err != nil {
			return err
		}
		body, err := bot.Request(durable.Ctx, "POST", "/messages", msg, accessToken, bot.UuidNewV4().String())

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
		adminIds := models.GetUserIdsByBotId(botInfo.ClientId)
		forwardDashboardMessage(&forwardMessagePropsType{
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

type forwardMessagePropsType struct {
	Category         string `json:"category"`
	CreatedAt        string `json:"created_at"`
	MessageId        string `json:"message_id"`
	Source           string `json:"source"`
	UserId           string `json:"user_id"`
	QuoteMessageId   string `json:"quote_message_id"`
	AdminId          string `json:"admin_id"`
	ConversationId   string `json:"conversation_id"`
	Data             string `json:"data"`
	RepresentativeId string `json:"representative_id"`
	Status           string `json:"status"`
	UpdatedAt        string `json:"updated_at"`
}

func forwardDashboardMessage(msg *forwardMessagePropsType, clientId, sessionId, privateKey, hash string, isAdmin bool, adminIds []string) {
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
		log.Println("attachment_id:  ", imageMessagePayload.AttachmentId)
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
			log.Println("forwardDashboardMessage 获取originMessage失败", msg.UserId)
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
			err := bot.PostMessages(durable.Ctx, sendMessages, clientId, sessionId, privateKey)
			if err != nil {
				log.Println("forwardDashboardMessage, 发送给其他管理员", err)
			}
		}
	} else if isAdmin {
		// 来自 客户端的 管理员 消息
		// 1.
		originMessage := models.GetOriginMessageById(clientId, msg.QuoteMessageId)
		if originMessage == nil || originMessage.RecipientId == "" {
			log.Println("forwardDashboardMessage 获取originMessage失败", msg.MessageId)
			return
		} else {
			userId = originMessage.RecipientId
		}
		status = "read"
	}

	userInfo, err := AddMessageUser(msg.UserId)
	if err != nil {
		log.Println(err)
	}

	// 发送给后台的管理员
	for _, chanel := range models.HashManagerMap[hash] {
		if chanel != nil {
			chanel <- models.RespMessage{
				ClientId:       clientId,
				IdentityNumber: userInfo.IdentityNumber,
				UserId:         userId,
				FullName:       userInfo.FullName,
				AvatarURL:      userInfo.AvatarURL,
				Data:           data,
				Category:       msg.Category,
				CreatedAt:      msg.CreatedAt,
				MessageId:      msg.MessageId,
				Source:         msg.Source,
				Status:         status,
			}
		}
	}
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

func getForwardParams() {

}
