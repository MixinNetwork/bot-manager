package models

import (
	"fmt"
	"github.com/liuzemei/bot-manager/db"
	"github.com/liuzemei/bot-manager/utils"
	"time"
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
	UserId         string `json:"user_id"`
}

type BotUser struct {
	ClientId  string `gorm:"column:client_id"`
	UserId    string `gorm:"column:user_id"`
	Status    string `gorm:"column:status"`
	BlockTime string `gorm:"column:block_time"`
	CreatedAt string `gorm:"column:created_at"`
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
  user_id      VARCHAR(36) NOT NULL,
  client_id    VARCHAR(36) NOT NULL,
  status			 VARCHAR DEFAULT '',
  block_time   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
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
	var botUser = BotUser{
		ClientId:  clientId,
		UserId:    u.UserId,
		BlockTime: utils.FormatTime(time.Now()),
		CreatedAt: utils.FormatTime(time.Now()),
	}
	db.Conn.Set("gorm:insert_option", "ON CONFLICT DO NOTHING").Create(&botUser)
}

func UpdateBotUserStatus(clientId, userId, status string) {
	if status == "normal" {
		status = ""
		db.Conn.Table("bot_users").Where("client_id=? AND user_id=?", clientId, userId).Update("status", status)
	} else {
		db.Conn.Table("bot_users").Where("client_id=? AND user_id=?", clientId, userId).Update(map[string]interface{}{"status": status, "block_time": time.Now()})
	}
}

func CheckUserStatus(clientId, userId string) bool {
	var botUser BotUser
	db.Conn.First(&botUser, "client_id=? AND user_id=?", clientId, userId)
	return botUser.Status == "block"
}

func GetTodayUserCount(clientId string) (count int) {
	t := utils.GetDate(0)
	db.Conn.Table("bot_users").Where("to_char(created_at, 'YYYY-MM-DD')=? AND client_id=?", t, clientId).Count(&count)
	return
}

func GetUserById(userId string) *UserBaseResp {
	var userInfo User
	db.Conn.First(&userInfo, "user_id=?", userId)
	return &UserBaseResp{
		FullName:       userInfo.FullName,
		IdentityNumber: userInfo.IdentityNumber,
		AvatarURL:      userInfo.AvatarURL,
		UserId:         userInfo.UserId,
	}
}

func GetUserByIds(userIds []string) []*UserBaseResp {
	userInfos := make([]*User, 0)
	db.Conn.Find(&userInfos, "user_id in (?)", userIds)
	resp := make([]*UserBaseResp, 0)
	for _, info := range userInfos {
		resp = append(resp, &UserBaseResp{
			FullName:       info.FullName,
			IdentityNumber: info.IdentityNumber,
			AvatarURL:      info.AvatarURL,
			UserId:         info.UserId,
		})
	}
	return resp
}

func GetBotUser(userId, clientId string) *UserBaseResp {
	var botUser UserBase
	db.Conn.Raw("select users.identity_number, users.avatar_url, users.full_name from bot_users left join users on bot_users.user_id=users.user_id where bot_users.user_id=? AND bot_users.client_id=?", userId, clientId).Scan(&botUser)
	return &UserBaseResp{
		FullName:       botUser.FullName,
		IdentityNumber: botUser.IdentityNumber,
		AvatarURL:      botUser.AvatarURL,
	}
}

type UserIdType struct {
	UserId string `gorm:"user_id"`
}

func GetBotUserListById(clientId string) []string {
	botUser := make([]*UserIdType, 0)
	db.Conn.Raw("select user_id from bot_users where client_id=?", clientId).Scan(&botUser)
	users := make([]string, 0)
	for _, userId := range botUser {
		users = append(users, userId.UserId)
	}
	return users
}

type BotUserType struct {
	IdentityNumber string `gorm:"column:identity_number"`
	AvatarURL      string `gorm:"column:avatar_url"`
	FullName       string `gorm:"column:full_name"`
	CreatedAt      string `gorm:"column:created_at"`
	ClientId       string `gorm:"column:client_id"`
	UserId         string `gorm:"column:user_id"`
	Status         string `gorm:"column:status"`
	BlockTime      string `gorm:"column:block_time"`
}
type BotUserTypeResp struct {
	IdentityNumber string `json:"identity_number"`
	AvatarURL      string `json:"avatar_url"`
	FullName       string `json:"full_name"`
	CreatedAt      string `json:"created_at"`
	UserId         string `json:"user_id"`
}
type BotUserBlackTypeResp struct {
	IdentityNumber string `json:"identity_number"`
	AvatarURL      string `json:"avatar_url"`
	FullName       string `json:"full_name"`
	CreatedAt      string `json:"created_at"`
	UserId         string `json:"user_id"`
	BlockTime      string `json:"block_time"`
}

func GetUsersByClientId(clientId, status string) interface{} {
	users := make([]*BotUserType, 0)
	if status == "normal" {
		status = ""
	}
	db.Conn.Raw("select users.user_id, users.identity_number, users.avatar_url, users.full_name, to_char(bot_users.created_at, 'YYYY/MM/DD HH24:MI:SS') as created_at, to_char(block_time, 'YYYY/MM/DD HH24:MI:SS') as block_time from bot_users left join users on bot_users.user_id=users.user_id where bot_users.client_id=? and bot_users.status=?", clientId, status).Scan(&users)
	if status == "" {
		userList := make([]*BotUserTypeResp, 0)
		for _, user := range users {
			userList = append(userList, &BotUserTypeResp{
				IdentityNumber: user.IdentityNumber,
				AvatarURL:      user.AvatarURL,
				FullName:       user.FullName,
				CreatedAt:      user.CreatedAt,
				UserId:         user.UserId,
			})
		}
		return userList
	} else {
		userList := make([]*BotUserBlackTypeResp, 0)
		for _, user := range users {
			userList = append(userList, &BotUserBlackTypeResp{
				IdentityNumber: user.IdentityNumber,
				AvatarURL:      user.AvatarURL,
				FullName:       user.FullName,
				CreatedAt:      user.CreatedAt,
				UserId:         user.UserId,
				BlockTime:      user.BlockTime,
			})
		}
		return userList
	}
}
