package externals

import (
	"context"
	"encoding/base64"
	"log"
	"strings"
	"time"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/astaxie/beego"
	"github.com/liuzemei/bot-manager/durable"
	"github.com/liuzemei/bot-manager/models"
	"github.com/liuzemei/bot-manager/utils"
)

type listener struct {
	Hash       string
	ClientId   string
	SessionId  string
	PrivateKey string
	ackMessage []*bot.ReceiptAcknowledgementRequest
}

var ignoreCategory = map[string]bool{
	"SYSTEM_CONVERSATION":     true,
	"SYSTEM_ACCOUNT_SNAPSHOT": true,
	"MESSAGE_RECALL":          true,
}

func (l listener) OnMessage(ctx context.Context, msg bot.MessageView, userId string) error {
	if ignoreCategory[msg.Category] {
		return nil
	}
	if models.CheckUserStatus(l.ClientId, msg.UserId) {
		return nil
	}
	data, _ := base64.StdEncoding.DecodeString(msg.Data)
	models.AddMessage(models.Message{
		ClientId:       l.ClientId,
		ConversationId: msg.ConversationId,
		UserId:         msg.UserId,
		MessageId:      msg.MessageId,
		Category:       msg.Category,
		Data:           string(data),
		Status:         strings.ToLower(msg.Status),
		Source:         msg.Source,
		CreatedAt:      utils.FormatTime(msg.CreatedAt),
	})
	models.HashMessengerMap[l.Hash] <- models.MessengerChannel{
		Message:    msg,
		ClientID:   l.ClientId,
		SessionID:  l.SessionId,
		PrivateKey: l.PrivateKey,
	}

	return nil
}

var ReconnectTimes = map[string]int{}

func StartWebSockets(clientId, sessionId, privateKey, hash string) error {
	for {
		client := bot.NewBlazeClient(clientId, sessionId, privateKey)
		var _listener = listener{
			ClientId:   clientId,
			SessionId:  sessionId,
			PrivateKey: privateKey,
			Hash:       hash,
		}
		err := client.Loop(durable.Ctx, _listener)
		if err != nil {
			if err.Error() == "websocket: bad handshake" {
				log.Println("StartWebSockets Err", err)
			} else if err.Error() == `{"status":500,"code":7000,"description":"Blaze server error."}` {
				if ReconnectTimes[clientId] >= 10 {
					return err
				}
				ReconnectTimes[clientId]++
				log.Println("Blaze server error", clientId, ReconnectTimes[clientId])
			} else if strings.Contains(err.Error(), "operation timed out") {
				log.Println("Blaze timed out")
			} else {
				log.Println("密码不对？", err.Error())
				return err
			}
		}
		time.Sleep(time.Second * 15)
	}
}

func SendBatchMessage(messages []*bot.MessageRequest, clientId, sessionId, privateKey string) error {
	if len(messages) <= 100 {
		err := bot.PostMessages(durable.Ctx, messages, clientId, sessionId, privateKey)
		if err != nil {
			return err
		}
	} else {
		currentMessage := make([]*bot.MessageRequest, 0)
		overMessage := messages
		for {
			if len(overMessage) > 100 {
				currentMessage = overMessage[:100]
				err := bot.PostMessages(durable.Ctx, currentMessage, clientId, sessionId, privateKey)
				if err != nil {
					log.Println("SendBatchMessage", err)
					continue
				}
				overMessage = overMessage[100:]
			} else {
				err := bot.PostMessages(durable.Ctx, overMessage, clientId, sessionId, privateKey)
				if err != nil {
					log.Println("SendBatchMessage2", err)
					continue
				}
				break
			}
		}
	}
	return nil
}

func SendText(botInfo *models.Bot, userId, text string) error {
	conversationId := bot.UniqueConversationId(botInfo.ClientId, userId)
	data := base64.StdEncoding.EncodeToString([]byte(text))
	err := bot.PostMessage(durable.Ctx, conversationId, userId, bot.UuidNewV4().String(), "PLAIN_TEXT", data, botInfo.ClientId, botInfo.SessionId, botInfo.PrivateKey)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	runmode := beego.AppConfig.String("runmode")
	if runmode == "prod" {
		bot.SetBlazeUri("blaze.mixin.one")
	}
}
