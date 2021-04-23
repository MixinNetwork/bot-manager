package models

import (
	"github.com/liuzemei/bot-manager/db"
	"github.com/liuzemei/bot-manager/utils"
	"time"
)

type AutoReplayMessage struct {
	ReplayId  string `gorm:"column:replay_id;type:varchar(36);" json:"replay_id,omitempty"`
	ClientId  string `gorm:"column:client_id" json:"client_id,omitempty"`
	Category  string `gorm:"column:category" json:"category,omitempty"`
	Data      string `gorm:"column:data" json:"data,omitempty"`
	Key       string `gorm:"column:key" json:"key,omitempty"`
	CreatedAt string `gorm:"column:created_at" json:"created_at,omitempty"`
}

func init() {
	db.RegisterMigration(`CREATE TABLE IF NOT EXISTS auto_replay_messages(
  replay_id     VARCHAR(36) NOT NULL,
  client_id     VARCHAR(36) NOT NULL,
  category      VARCHAR(36) NOT NULL,
  data          TEXT NOT NULL,
  key           VARCHAR NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	PRIMARY KEY(replay_id, key)
);`)

	db.RegisterModel(&AutoReplayMessage{})
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
		db.Conn.Model(&autoReplayMessage).Where("client_id=? AND replay_id=?", clientId, replayId).Update(map[string]interface{}{"data": data, "category": category})
	}
}

func GetAutoReplayMessage(clientId string) []*AutoReplayMessage {
	var autoReplayMessages []*AutoReplayMessage
	db.Conn.Order("created_at ASC").Find(&autoReplayMessages, "client_id=?", clientId)
	return autoReplayMessages
}
func DeleteAutoReplayMessage(replayId string) {
	db.Conn.Delete(&AutoReplayMessage{}, "replay_id=?", replayId)
}

func GetAutoReplayMessageByKey(clientId, key string) (string, string) {
	var autoReplayMessage AutoReplayMessage
	db.Conn.Order("created_at DESC").First(&autoReplayMessage, "client_id=? AND key=?", clientId, key)
	return autoReplayMessage.Data, autoReplayMessage.Category
}
