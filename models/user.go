package models

import (
	"fmt"
	"github.com/liuzemei/bot-manager/db"
	"github.com/liuzemei/bot-manager/utils"
)

type User struct {
	UserId         string `gorm:"column:user_id"`
	FullName       string `gorm:"column:full_name"`
	IdentityNumber string `gorm:"column:identity_number"`
	AvatarURL      string `gorm:"column:avatar_url"`
	AccessToken    string `gorm:"column:access_token"`
	CreatedAt      string `gorm:"column:created_at"`
}

type UserBase struct {
	FullName       string `gorm:"column:full_name"`
	IdentityNumber string `gorm:"column:identity_number"`
	AvatarURL      string `gorm:"column:avatar_url"`
}

type UserBaseResp struct {
	FullName       string `json:"full_name"`
	IdentityNumber string `json:"identity_number"`
	AvatarURL      string `json:"avatar_url"`
}

type BotUser struct {
	ClientId string `gorm:"column:client_id"`
	UserId   string `gorm:"column:user_id"`
}

func init() {
	db.RegisterModel(&User{})
	db.RegisterModel(&BotUser{})
	db.RegisterMigration(`CREATE TABLE IF NOT EXISTS users (
  user_id             VARCHAR(36) NOT NULL PRIMARY KEY,
  full_name           VARCHAR(1024) NOT NULL,
  identity_number     VARCHAR(11) NOT NULL UNIQUE,
  avatar_url          VARCHAR(1024) NOT NULL,
  access_token        VARCHAR(512) NOT NULL DEFAULT '',
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);`)
	db.RegisterMigration(`CREATE TABLE IF NOT EXISTS bot_users (
  user_id             VARCHAR(36) NOT NULL,
  client_id           VARCHAR(36) NOT NULL,
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	PRIMARY KEY(user_id, client_id)
);`)
}

func AddUser(u *User) {
	var _u User
	var updateStr string
	db.Conn.First(&_u, "user_id=?", u.UserId)

	if _u.UserId != "" {
		if u.AccessToken != "" {
			updateStr = fmt.Sprintf(
				"ON CONFLICT(%s) DO UPDATE SET full_name='%s', identity_number='%s', avatar_url='%s', access_token='%s'",
				"user_id", u.FullName, u.IdentityNumber, u.AvatarURL, u.AccessToken)
		} else {
			updateStr = fmt.Sprintf(
				"ON CONFLICT(%s) DO UPDATE SET full_name='%s', identity_number='%s', avatar_url='%s'",
				"user_id", u.FullName, u.IdentityNumber, u.AvatarURL)
		}
		db.Conn.Set("gorm:insert_option", updateStr).Create(&u)
	} else {
		db.Conn.Create(&u)
	}
}


func AddBotUser(u *User, clientId string) {
	AddUser(u)
	var botUser = BotUser{
		ClientId: clientId,
		UserId:   u.UserId,
	}
	db.Conn.Set("gorm:insert_option", "ON CONFLICT DO NOTHING").Create(&botUser)
}

func DeleteBotUser(userID, clientId string) {
	db.Conn.Delete(BotUser{}, "user_id=? AND client_id=?", userID, clientId)
}

func GetTodayUserCount(clientId string) (count int) {
	t := utils.GetDate(0)
	db.Conn.Debug().Table("bot_users").Where("to_char(created_at, 'YYYY-MM-DD')=? AND client_id=?", t, clientId).Count(&count)
	return
}
//
//func GetUserById(userID string) *UserBaseResp {
//	var userInfo UserBaseResp
//	db.Conn.First(&userInfo, "user_id=?", userID)
//	return &userInfo
//}
