package models

import (
	"time"

	"github.com/MixinNetwork/bot-manager/db"
	"github.com/MixinNetwork/bot-manager/utils"
)

type Broadcast struct {
	ClientId  string `gorm:"column:client_id;type:varchar(36);not null;primaryKey"`
	UserId    string `gorm:"column:user_id;type:varchar(36);not null"`
	MessageId string `gorm:"column:message_id;type:varchar(36);not null;primaryKey"`
	Category  string `gorm:"column:category;type:varchar(36);not null"`
	Data      string `gorm:"column:data;type:text;not null"`
	CreatedAt string `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()"`
}

type BroadcastTmp struct {
	ClientId        string `gorm:"column:client_id;type:varchar(36);not null"`
	MessageId       string `gorm:"column:message_id;type:varchar(36);not null"`
	OriginMessageId string `gorm:"column:origin_message_id;type:varchar(36);not null;index"`
	UserId          string `gorm:"column:user_id;type:varchar(36);not null"`
	ConversationId  string `gorm:"column:conversation_id;type:varchar(36);not null"`
}

func init() {
	db.RegisterMigration(`
CREATE TABLE IF NOT EXISTS broadcasts(
  client_id     VARCHAR(36) NOT NULL,
	user_id       VARCHAR(36) NOT NULL,
  message_id    VARCHAR(36) NOT NULL,
  category      VARCHAR(36) NOT NULL,
  data          TEXT NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY(client_id, message_id)
);

CREATE TABLE IF NOT EXISTS broadcast_tmps(
  client_id              VARCHAR(36) NOT NULL,
  message_id             VARCHAR(36) NOT NULL,
  user_id                VARCHAR(36) NOT NULL,
  conversation_id        VARCHAR(36) NOT NULL,
  origin_message_id      VARCHAR(36) NOT NULL
);
CREATE INDEX IF NOT EXISTS broadcast_tmp_origin_message_id_idx ON broadcast_tmps(origin_message_id);
`)
	db.RegisterModel(&BroadcastTmp{})
	db.RegisterModel(&Broadcast{})
}

func AddBroadcast(clientId, userId, messageId, category, data string) {
	db.Conn.Create(&Broadcast{
		ClientId:  clientId,
		UserId:    userId,
		MessageId: messageId,
		Category:  category,
		Data:      data,
		CreatedAt: utils.FormatTime(time.Now()),
	})
}

type BroadcastResp struct {
	MessageId string `json:"message_id"`
	Category  string `json:"category"`
	Data      string `json:"data"`
	CreatedAt string `json:"created_at"`
	FullName  string `json:"full_name"`
}

func GetBroadcast(clientId string) []*BroadcastResp {
	broadcast := make([]*struct {
		FullName  string `gorm:"column:full_name"`
		MessageId string `gorm:"column:message_id"`
		Category  string `gorm:"column:category"`
		Data      string `gorm:"column:data"`
		CreatedAt string `gorm:"column:created_at"`
	}, 0)
	db.Conn.Raw("SELECT message_id,category,data,to_char(broadcasts.created_at, 'YYYY/MM/DD HH24:MI:SS') as created_at, full_name  FROM broadcasts LEFT JOIN users ON users.user_id=broadcasts.user_id WHERE client_id=? ORDER BY created_at DESC", clientId).Scan(&broadcast)
	target := make([]*BroadcastResp, 0)
	for _, b := range broadcast {
		target = append(target, &BroadcastResp{
			MessageId: b.MessageId,
			Category:  b.Category,
			Data:      b.Data,
			FullName:  b.FullName,
			CreatedAt: b.CreatedAt,
		})
	}
	return target
}

func DeleteBroadcast(clientId, messageId string) {
	db.Conn.Delete(&Broadcast{}, "client_id=? AND message_id=?", clientId, messageId)
}

func AddBroadcastTmpMessage(clientId, messageId, originMessageId, userId, conversationId string) {
	db.Conn.Create(&BroadcastTmp{
		ClientId:        clientId,
		MessageId:       messageId,
		OriginMessageId: originMessageId,
		ConversationId:  conversationId,
		UserId:          userId,
	})
}

func GetBroadcastTmpMessage(clientId, originMessageId string) []*BroadcastTmp {
	broadcastTmp := make([]*BroadcastTmp, 0)
	db.Conn.Find(&broadcastTmp, "client_id=? AND origin_message_id=?", clientId, originMessageId)
	return broadcastTmp
}

func DeleteBroadcastTmp(clientId, originMessageId string) {
	db.Conn.Delete(&BroadcastTmp{}, "client_id=? AND origin_message_id=?", clientId, originMessageId)
}
