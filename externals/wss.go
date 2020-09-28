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

func (l listener) OnMessage(ctx context.Context, msg bot.MessageView, userId string) error {

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
	client := bot.NewBlazeClient(clientId, sessionId, privateKey)
	var _listener = listener{
		ClientId:   clientId,
		SessionId:  sessionId,
		PrivateKey: privateKey,
		Hash:       hash,
	}
	for {
		err := client.Loop(durable.Ctx, _listener)
		if err != nil {
			if err.Error() == "websocket: bad handshake" {
				log.Println("StartWebSockets Err", err)
			} else {
				log.Println("密码不对？", err)
				return err
			}
		}
	}
}
