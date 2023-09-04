package models

import (
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/bot-manager/db"
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
	ClientId       string `gorm:"column:client_id;type:varchar(36);not null;primaryKey" json:"client_id,omitempty"`
	UserId         string `gorm:"column:user_id;type:varchar(36);not null" json:"user_id,omitempty"`
	ConversationId string `gorm:"column:conversation_id;type:varchar(36);not null" json:"conversation_id,omitempty"`
	MessageId      string `gorm:"column:message_id;type:varchar(36);not null;primaryKey" json:"message_id,omitempty"`
	Category       string `gorm:"column:category;type:varchar(36)" json:"category,omitempty"`
	Data           string `gorm:"column:data;type:text" json:"data,omitempty"`
	Status         string `gorm:"column:status;type:varchar(36);not null;" json:"status,omitempty"`
	Source         string `gorm:"column:source;type:varchar(36);not null;" json:"source,omitempty"`
	CreatedAt      string `gorm:"column:created_at;type:timestamp with time zone;not null;default now();" json:"created_at,omitempty"`

	IdentityNumber   string `json:"identity_number,omitempty"`
	QuoteMessageId   string `json:"quote_message_id,omitempty"`
	AdminId          string `json:"admin_id,omitempty"`
	RepresentativeId string `json:"representative_id,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}
type RespMessage struct {
	Message

	FullName    string `json:"full_name,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	RecipientId string `json:"recipient_id"`
}

type ForwardMessage struct {
	ClientId        string `gorm:"column:client_id;type:varchar(36);not null;primaryKey" json:"client_id,omitempty"`
	MessageId       string `gorm:"column:message_id;type:varchar(36);not null;primaryKey" json:"message_id,omitempty"`
	AdminId         string `gorm:"column:admin_id;type:varchar(36);not null" json:"admin_id,omitempty"`
	RecipientId     string `gorm:"column:recipient_id;type:varchar(36);not null" json:"recipient_id,omitempty"`
	OriginMessageId string `gorm:"column:origin_message_id;type:varchar(36);not null" json:"origin_message_id,omitempty"`
	ConversationId  string `gorm:"column:conversation_id;type:varchar(36);not null" json:"conversation_id,omitempty"`
	AdminMessageId  string `gorm:"column:admin_message_id;type:varchar(36);not null" json:"admin_message_id,omitempty"`
	CreatedAt       string `gorm:"column:created_at;type:varchar(36);not null" json:"created_at,omitempty"`
}

func init() {
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
	var messages []RespMessage
	db.Conn.Table("messages").Select("users.identity_number, users.full_name, users.avatar_url, messages.*").Joins("left join users on messages.user_id=users.user_id").Where("client_id IN (?) AND messages.created_at-?>interval '0 day' ", clientIds, date).Order("created_at ASC").Find(&messages)
	return messages
}

func GetAllMessagesByBotId(clientId string) []RespMessage {
	var msgs []RespMessage
	db.Conn.Table("messages").Select("users.identity_number, users.full_name, users.avatar_url, messages.*").Joins("left join users on messages.user_id=users.user_id").Where("client_id=?", clientId).Order("created_at ASC").Find(&msgs)
	return msgs
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

func UpdateClientMessageById(clientId string, user *User, msg *Message, status string) {
	db.Conn.Exec(`UPDATE messages SET user_id=$3, status=$4 WHERE client_id=$1 AND message_id=$2`, clientId, msg.MessageId, user.UserId, status)
}
