package models

import (
	"github.com/liuzemei/bot-manager/db"
	"github.com/liuzemei/bot-manager/utils"
)

type DailyData struct {
	ClientId string `gorm:"column:client_id"`
	Date     string `gorm:"column:date"`
	Users    int    `gorm:"column:users"`
	Messages int    `gorm:"column:messages"`
}

type RespDailyData struct {
	Date     string `json:"date"`
	Users    int    `json:"users"`
	Messages int    `json:"messages"`
}

type RespGetData struct {
	List  []RespDailyData `json:"list"`
	Today RespDailyData   `json:"today"`
}

func init() {

	//UserList = make(map[string]*User)
	//u := User{"user_11111", "astaxie", "11111", Profile{"male", 20, "Singapore", "astaxie@gmail.com"}}
	//UserList["user_11111"] = &u

	db.RegisterModel(&DailyData{})
	db.RegisterMigration(`CREATE TABLE IF NOT EXISTS daily_data (
  client_id  VARCHAR(36) NOT NULL,
	date       DATE NOT NULL,
	users      INTEGER NOT NULL DEFAULT 0,
	messages   INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY(client_id, date)
);`)
}

func AddDailyData(clientId string, users int, messages int) {
	var dailyData = DailyData{
		ClientId: clientId,
		Date:     "",
		Users:    users,
		Messages: messages,
	}
	db.Conn.Create(&dailyData)
}

type count = struct {
	Count int `gorm:"column:count"`
}

func GetDailyData(clientID string) RespGetData {
	var dailyList []DailyData
	var resp RespGetData
	db.Conn.Table("daily_data").Select("date, users, messages").Where("client_id=?", clientID).Scan(&dailyList)
	for _, data := range dailyList {
		resp.List = append(resp.List, RespDailyData{
			Date:     data.Date,
			Users:    data.Users,
			Messages: data.Messages,
		})
	}
	var messageCount count
	db.Conn.Debug().Raw("SELECT count(1) FROM messages WHERE client_id=? AND created_at - CURRENT_DATE > interval ' 0 day'", clientID).Scan(&messageCount)
	var userCount count
	db.Conn.Raw("SELECT count(1) FROM bot_users WHERE client_id=? AND created_at - CURRENT_DATE > interval ' 0 day'", clientID).Scan(&userCount)
	resp.Today.Messages = messageCount.Count
	resp.Today.Users = userCount.Count
	resp.Today.Date = utils.GetDate(0)
	return resp
}
