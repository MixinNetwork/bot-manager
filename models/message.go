package models

import (
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/liuzemei/bot-manager/db"
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
type ForwardMessagePropsType struct {
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
	db.RegisterModel(&Message{})
	db.RegisterModel(&ForwardMessage{})
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

func GetAllMessagesByUserId(userId string, date string) []RespMessage {
	clientIds := GetBotIdsByUserId(userId)
	var messages []RespDbMessage
	db.Conn.Table("messages").Select("users.identity_number, users.full_name, users.avatar_url, messages.*").Joins("left join users on messages.user_id=users.user_id").Where("client_id IN (?) AND messages.created_at-?>interval '0 day' ", clientIds, date).Order("created_at ASC").Find(&messages)
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
	db.Conn.Select("message_id, conversation_id, origin_message_id, recipient_id").Where("client_id=? AND message_id=?", clientId, messageId).First(&msg)
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
	db.Conn.Select("origin_message_id").Where("client_id=? AND recipient_id=?", clientId, recipientId).Order("created_at DESC").First(&msg)
	if msg.OriginMessageId == "" {
		return nil
	}
	return &msg
}

func UpdateClientMessageById(clientId string, user *UserBaseResp, msg *ForwardMessagePropsType, status string) {
	db.Conn.Exec(`UPDATE messages SET user_id=$3, status=$4 WHERE client_id=$1 AND message_id=$2`, clientId, msg.MessageId, user.UserId, status)
}
