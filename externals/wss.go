package externals

import (
	"context"
	"encoding/base64"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/liuzemei/bot-manager/durable"
	"github.com/liuzemei/bot-manager/models"
	"github.com/liuzemei/bot-manager/utils"
	"log"
	"strings"
)

type listener struct {
	Hash       string
	ClientId   string
	SessionId  string
	PrivateKey string
	ackMessage []*bot.ReceiptAcknowledgementRequest
}

var ignoreCategory = map[string]bool{
	"SYSTEM_CONVERSATION": true,
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
				log.Println(clientId)
				log.Println("Blaze server error")
			} else if strings.Contains(err.Error(), "operation timed out") {
				log.Println("Blaze timed out")
			} else {
				log.Println("密码不对？", err.Error())
				return nil
			}
		}
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
		overMessage := make([]*bot.MessageRequest, 0)
		for {
			if len(overMessage) > 100 {
				currentMessage = overMessage[:100]
				err := bot.PostMessages(durable.Ctx, currentMessage, clientId, sessionId, privateKey)
				if err != nil {
					log.Println(err)
					continue
				}
				overMessage = overMessage[100:]
			} else {
				err := bot.PostMessages(durable.Ctx, currentMessage, clientId, sessionId, privateKey)
				if err != nil {
					log.Println(err)
					continue
				}
				break
			}
		}
	}
	return nil
}

func init() {
	bot.SetBlazeUri("blaze.mixin.one")
}
