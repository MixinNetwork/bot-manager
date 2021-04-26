package models

import (
	"github.com/MixinNetwork/bot-manager/db"
	"github.com/MixinNetwork/bot-manager/utils"
)

type DailyData struct {
	ClientId string `gorm:"column:client_id;type:varchar(36);not null;primaryKey" json:"client_id"`
	Date     string `gorm:"column:date;type:date;not null;primaryKey" json:"date"`
	Users    int    `gorm:"column:users;type:integer;not null;default 0" json:"users"`
	Messages int    `gorm:"column:messages;type:integer;not null;default 0" json:"messages"`
}

type RespGetData struct {
	List  []*DailyData `json:"list"`
	Today *DailyData   `json:"today"`
}

func init() {
	db.RegisterModel(&DailyData{})
	db.RegisterMigration(`CREATE TABLE IF NOT EXISTS daily_data (
  client_id  VARCHAR(36) NOT NULL,
	date       DATE NOT NULL,
	users      INTEGER NOT NULL DEFAULT 0,
	messages   INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY(client_id, date)
);`)
}

func AddDailyData(clientId string, users int, messages int, date string) {
	db.Conn.Create(&DailyData{
		ClientId: clientId,
		Date:     date,
		Users:    users,
		Messages: messages,
	})
}

type count = struct {
	Count int `gorm:"column:count"`
}

func GetDailyData(clientID string) RespGetData {
	var dailyList []DailyData
	var resp RespGetData
	resp.List = make([]*DailyData, 0)
	db.Conn.Table("daily_data").Select("to_char(date, 'YYYY-MM-DD') as date, users, messages").Where("client_id=?", clientID).Scan(&dailyList)
	for _, data := range dailyList {
		resp.List = append(resp.List, &DailyData{
			Date:     data.Date,
			Users:    data.Users,
			Messages: data.Messages,
		})
	}
	resp.Today = GetTodayData(clientID)
	return resp
}

func GetTodayData(clientID string) *DailyData {
	var messageCount count
	db.Conn.Raw("SELECT count(1) FROM messages WHERE client_id=? AND created_at - CURRENT_DATE > interval ' 0 day'", clientID).Scan(&messageCount)
	var userCount count
	db.Conn.Raw("SELECT count(1) FROM bot_users WHERE client_id=? AND created_at - CURRENT_DATE > interval ' 0 day'", clientID).Scan(&userCount)
	return &DailyData{
		Date:     utils.GetDate(0),
		Users:    userCount.Count,
		Messages: messageCount.Count,
	}
}

func SaveTodayData(num int) error {
	users := make([]*struct {
		ClientId string `gorm:"client_id"`
		Users    int    `gorm:"users"`
	}, 0)
	db.Conn.Raw(`SELECT client_id, count(1) as users FROM bot_users WHERE to_char(created_at, 'YYYY-MM-DD')=? GROUP BY client_id`, utils.GetDate(num)).Scan(&users)
	messages := make([]*struct {
		ClientId string `gorm:"client_id"`
		Messages int    `gorm:"messages"`
	}, 0)
	db.Conn.Raw(`SELECT client_id, count(1) as messages FROM messages WHERE to_char(created_at, 'YYYY-MM-DD')=? GROUP BY client_id`, utils.GetDate(num)).Scan(&messages)
	target := make([]*DailyData, 0)
	hasHandle := map[string]bool{}
	for _, user := range users {
		hasHandle[user.ClientId] = true
		_m := 0
		for _, message := range messages {
			if message.ClientId == user.ClientId {
				_m = message.Messages
				break
			}
		}
		target = append(target, &DailyData{
			ClientId: user.ClientId,
			Date:     utils.GetDate(-1),
			Users:    user.Users,
			Messages: _m,
		})
	}
	for _, message := range messages {
		if hasHandle[message.ClientId] {
			continue
		}
		_u := 0
		for _, user := range users {
			if user.ClientId == message.ClientId {
				_u = user.Users
			}
		}
		target = append(target, &DailyData{
			ClientId: message.ClientId,
			Date:     utils.GetDate(-1),
			Users:    _u,
			Messages: message.Messages,
		})
	}
	for _, data := range target {
		AddDailyData(data.ClientId, data.Users, data.Messages, data.Date)
	}
	return nil
}
