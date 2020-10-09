package models

import (
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/liuzemei/bot-manager/db"
	"github.com/liuzemei/bot-manager/utils"
	"time"
)

type MessengerChannel struct {
	Message    bot.MessageView
	ClientID   string
	SessionID  string
	PrivateKey string
}

var HashMessengerMap = make(map[string]chan MessengerChannel)
var HashManagerMap = make(map[string]map[string]chan RespMessage)

type MessageListener struct {
	Hash       string
	ClientId   string
	SessionId  string
	PrivateKey string
	ackMessage []*bot.ReceiptAcknowledgementRequest
}

type Message struct {
	ClientId       string `gorm:"column:client_id"`
	UserId         string `gorm:"column:user_id"`
	ConversationId string `gorm:"column:conversation_id"`
	MessageId      string `gorm:"column:message_id"`
	Category       string `gorm:"column:category"`
	Data           string `gorm:"column:data"`
	Status         string `gorm:"column:status"`
	Source         string `gorm:"column:source"`
	CreatedAt      string `gorm:"column:created_at"`
}
type RespDbMessage struct {
	ClientId       string `gorm:"column:client_id"`
	UserId         string `gorm:"column:user_id"`
	IdentityNumber string `gorm:"column:identity_number"`
	FullName       string `gorm:"column:full_name"`
	AvatarURL      string `gorm:"column:avatar_url"`
	ConversationId string `gorm:"column:conversation_id"`
	MessageId      string `gorm:"column:message_id"`
	Category       string `gorm:"column:category"`
	Data           string `gorm:"column:data"`
	Status         string `gorm:"column:status"`
	Source         string `gorm:"column:source"`
	CreatedAt      string `gorm:"column:created_at"`
}

type RespMessage struct {
	ClientId       string      `json:"client_id"`
	UserId         string      `json:"user_id"`
	RecipientId    string      `json:"recipient_id"`
	IdentityNumber string      `json:"identity_number"`
	FullName       string      `json:"full_name"`
	AvatarURL      string      `json:"avatar_url"`
	ConversationID string      `json:"conversation_id"`
	MessageId      string      `json:"message_id"`
	Category       string      `json:"category"`
	Data           interface{} `json:"data"`
	Status         string      `json:"status"`
	Source         string      `json:"source"`
	CreatedAt      string      `json:"created_at"`
}

type ForwardMessage struct {
	ClientId        string `gorm:"column:client_id"`
	MessageId       string `gorm:"column:message_id"`
	AdminId         string `gorm:"column:admin_id"`
	RecipientId     string `gorm:"column:recipient_id"`
	OriginMessageId string `gorm:"column:origin_message_id"`
	ConversationId  string `gorm:"column:conversation_id"`
	AdminMessageId  string `gorm:"column:admin_message_id"`
	CreatedAt       string `gorm:"column:created_at"`
}

type RespForwardMessage struct {
	ClientId        string `json:"client_id"`
	MessageId       string `json:"message_id"`
	AdminId         string `json:"admin_id"`
	RecipientId     string `json:"recipient_id"`
	OriginMessageId string `json:"origin_message_id"`
	ConversationId  string `json:"conversation_id"`
	AdminMessageId  string `json:"admin_message_id"`
	CreatedAt       string `json:"created_at"`
}

type AutoReplayMessage struct {
	ReplayId  string `gorm:"column:replay_id"`
	ClientId  string `gorm:"column:client_id"`
	Category  string `gorm:"column:category"`
	Data      string `gorm:"column:data"`
	Key       string `gorm:"column:key"`
	CreatedAt string `gorm:"column:created_at"`
}

type RespReplayMessage struct {
	ReplayId  string `json:"replay_id"`
	ClientId  string `json:"client_id"`
	Category  string `json:"category"`
	Data      string `json:"data"`
	Key       string `json:"key"`
	CreatedAt string `json:"created_at"`
}

func init() {
	db.RegisterMigration(`CREATE TABLE IF NOT EXISTS forward_messages (
  client_id              VARCHAR(36) NOT NULL,
  message_id             VARCHAR(36) NOT NULL,
  admin_id               VARCHAR(36) NOT NULL,
  recipient_id           VARCHAR(36) NOT NULL,
  origin_message_id      VARCHAR(36) NOT NULL,
  conversation_id        VARCHAR(36) NOT NULL,
  admin_message_id       VARCHAR(36),
  created_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY(client_id, message_id)
);`)
	db.RegisterMigration(`CREATE TABLE IF NOT EXISTS messages (
  client_id           VARCHAR(36) NOT NULL,
  user_id             VARCHAR(36) NOT NULL,
  conversation_id     VARCHAR(36) NOT NULL,
  message_id          VARCHAR(36) NOT NULL,
  category            VARCHAR(36),
  data                TEXT,
  status              VARCHAR(36) NOT NULL,
  source              VARCHAR(36) NOT NULL,
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY(client_id, message_id)
);`)
	db.RegisterMigration(`CREATE TABLE IF NOT EXISTS auto_replay_messages(
  replay_id     VARCHAR(36) NOT NULL,
  client_id     VARCHAR(36) NOT NULL,
  category      VARCHAR(36) NOT NULL,
  data          TEXT NOT NULL,
  key           VARCHAR NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	PRIMARY KEY(replay_id, key)
);`)
	db.RegisterModel(&Message{})
	db.RegisterModel(&ForwardMessage{})
	db.RegisterModel(&AutoReplayMessage{})
}
func AddMessage(message Message) {
	db.Conn.Set("gorm:insert_option", "ON CONFLICT DO NOTHING").Create(&message)
}

func AddForwardMessage(message ForwardMessage) {
	db.Conn.Create(&message)
}

func UpdateMessage(messageId, status string) {
	db.Conn.Table("messages").Where("message_id=?", messageId).Update("status", status)
}

func AddOrUpdateAutoReplayMessage(replayId, key, clientId, category, data string) {
	var autoReplayMessage AutoReplayMessage
	db.Conn.Order("created_at DESC").First(&autoReplayMessage, "client_id=? AND key=?", clientId, key)
	if autoReplayMessage.Category == "" {
		db.Conn.Create(&AutoReplayMessage{
			ReplayId:  replayId,
			ClientId:  clientId,
			Category:  category,
			Data:      data,
			Key:       key,
			CreatedAt: utils.FormatTime(time.Now()),
		})
	} else {
		db.Conn.Model(&autoReplayMessage).Update(map[string]interface{}{"data": data, "category": category})
	}
}

func GetAutoReplayMessage(clientId string) []*RespReplayMessage {
	var autoReplayMessages []*AutoReplayMessage
	db.Conn.Order("created_at ASC").Find(&autoReplayMessages, "client_id=?", clientId)
	resp := make([]*RespReplayMessage, 0)
	for _, message := range autoReplayMessages {
		resp = append(resp, &RespReplayMessage{
			ReplayId:  message.ReplayId,
			ClientId:  message.ClientId,
			Category:  message.Category,
			Data:      message.Data,
			Key:       message.Key,
			CreatedAt: message.CreatedAt,
		})
	}
	return resp
}
func DeleteAutoReplayMessage(replayId string) {
	db.Conn.Delete(&AutoReplayMessage{}, "replay_id=?", replayId)
}

func GetAutoReplayMessageByKey(clientId, key string) (string, string) {
	var autoReplayMessage AutoReplayMessage
	db.Conn.Order("created_at DESC").First(&autoReplayMessage, "client_id=? AND key=?", clientId, key)
	return autoReplayMessage.Data, autoReplayMessage.Category
}

func GetAllMessagesByUserId(userId string, date string) []RespMessage {
	clientIds := GetBotIdsByUserId(userId)
	var messages []RespDbMessage
	db.Conn.Debug().Table("messages").Select("users.identity_number, users.full_name, users.avatar_url, messages.*").Joins("left join users on messages.user_id=users.user_id").Where("client_id IN (?) AND messages.created_at-?>interval '0 day' ", clientIds, date).Order("created_at ASC").Find(&messages)
	var resp []RespMessage
	for _, message := range messages {
		resp = append(resp, RespMessage{
			ClientId:       message.ClientId,
			UserId:         message.UserId,
			IdentityNumber: message.IdentityNumber,
			FullName:       message.FullName,
			AvatarURL:      message.AvatarURL,
			ConversationID: message.ConversationId,
			MessageId:      message.MessageId,
			Category:       message.Category,
			Data:           message.Data,
			Status:         message.Status,
			Source:         message.Source,
			CreatedAt:      message.CreatedAt,
		})
	}
	return resp
}

func GetOriginMessageById(clientId, messageId string) *ForwardMessage {
	var msg ForwardMessage
	db.Conn.Debug().Select("message_id, conversation_id, origin_message_id, recipient_id").Where("client_id=? AND message_id=?", clientId, messageId).First(&msg)
	if msg.OriginMessageId == "" {
		return nil
	}
	return &msg
}

func GetForwardMessagesByOrigin(clientId, originMessageId string) map[string]ForwardMessage {
	var resp []ForwardMessage
	db.Conn.Select("admin_id, message_id").Where("client_id=? AND origin_message_id=?", clientId, originMessageId).Find(&resp)
	if len(resp) == 0 {
		return nil
	}
	var forwardMessagesList = map[string]ForwardMessage{}
	for _, message := range resp {
		forwardMessagesList[message.AdminId] = message
	}
	return forwardMessagesList
}

func GetLastMessageByRecipientId(clientId, recipientId string) *ForwardMessage {
	var msg ForwardMessage
	db.Conn.Debug().Select("origin_message_id").Where("client_id=? AND recipient_id=?", clientId, recipientId).Order("created_at DESC").First(&msg)
	if msg.OriginMessageId == "" {
		return nil
	}
	return &msg
}

func GetMessagesByOriginId(messageId string) *RespMessage {
	var message *RespMessage
	db.Conn.Raw("SELECT message_id, recipient_id, conversation_id, admin_id FROM message WHERE origin_message_id=?", messageId).Scan(message)
	return message
}

func transferForwardMessageData(messages []ForwardMessage) []RespForwardMessage {
	var resp []RespForwardMessage
	for _, message := range messages {
		resp = append(resp, RespForwardMessage{
			ClientId:        message.ClientId,
			MessageId:       message.MessageId,
			AdminId:         message.AdminId,
			RecipientId:     message.RecipientId,
			OriginMessageId: message.OriginMessageId,
			ConversationId:  message.ConversationId,
			AdminMessageId:  message.AdminMessageId,
			CreatedAt:       message.CreatedAt,
		})
	}
	return resp
}
