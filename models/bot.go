package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/liuzemei/bot-manager/db"
	"strings"
)

var ChangeBotWss = make(map[string]chan string)

type UserBot struct {
	UserId     string `gorm:"column:user_id"`
	ClientId   string `gorm:"column:client_id"`
	SessionId  string `gorm:"column:session_id"`
	PrivateKey string `gorm:"column:private_key"`
	Hash       string `gorm:"column:hash"`
}

type BotInfoRes struct {
	FullName       string `json:"full_name"`
	IdentityNumber string `json:"identity_number"`
	AvatarURL      string `json:"avatar_url"`
	ClientId       string `json:"client_id"`
}

type Bot struct {
	ClientId       string `gorm:"column:client_id"`
	SessionId      string `gorm:"column:session_id"`
	PrivateKey     string `gorm:"column:private_key"`
	FullName       string `gorm:"column:full_name"`
	IdentityNumber string `gorm:"column:identity_number"`
	AvatarURL      string `gorm:"column:avatar_url"`
	Hash           string `gorm:"column:hash"`
}

func init() {
	db.RegisterModel(UserBot{})
	db.RegisterModel(Bot{})

	db.RegisterMigration(`
CREATE TABLE IF NOT EXISTS user_bots (
  user_id      VARCHAR(36) NOT NULL,
  client_id    VARCHAR(36),
  session_id   VARCHAR(36),
  private_key  VARCHAR,
  hash         VARCHAR NOT NULL,
  PRIMARY KEY(user_id, client_id)
);`)

	db.RegisterMigration(`
CREATE TABLE IF NOT EXISTS bots (
  client_id           VARCHAR(36) NOT NULL PRIMARY KEY,
  session_id          VARCHAR(36) NOT NULL,
  private_key         VARCHAR NOT NULL,
  full_name           VARCHAR(1024) NOT NULL,
  identity_number     VARCHAR(11) NOT NULL UNIQUE,
  avatar_url          VARCHAR(1024) NOT NULL,
  hash                VARCHAR NOT NULL UNIQUE
);`)
}

func AddOrUpdateUserBotItem(userId, clientId, sessionId, privateKey string) {
	hash := Sha256Hash(clientId, sessionId, privateKey)
	var userBot = UserBot{
		UserId:     userId,
		ClientId:   clientId,
		SessionId:  sessionId,
		PrivateKey: privateKey,
		Hash:       hash,
	}
	updateStr := fmt.Sprintf("ON CONFLICT(%s, %s) DO UPDATE SET session_id='%s', private_key='%s', hash='%s'", "user_id", "client_id", sessionId, privateKey, hash)
	db.Conn.Set("gorm:insert_option", updateStr).Create(&userBot)
}

func GetBotListByUserId(userId string) []UserBase {
	var userBotList []UserBot
	db.Conn.Table("user_bots").Select("hash").Where("user_id=?", userId).Scan(&userBotList)
	var hashList []string
	for _, k := range userBotList {
		hashList = append(hashList, k.Hash)
	}
	botList := GetBotByHash(hashList)
	return botList
}

func GetAdminIdsByBotId(clientId string) []string {
	var userBots []UserBot
	db.Conn.Table("user_bots").Select("user_id").Where("client_id=?", clientId).Scan(&userBots)
	var userIds []string
	for _, bot := range userBots {
		userIds = append(userIds, bot.UserId)
	}
	return userIds
}

func GetBotIdsByUserId(userId string) []string {
	var userBots []UserBot
	db.Conn.Table("user_bots").Select("client_id").Where("user_id=?", userId).Scan(&userBots)
	var clientIds []string
	for _, bot := range userBots {
		clientIds = append(clientIds, bot.ClientId)
	}
	return clientIds
}

func GetUserBotHashByUserId(userId string) []string {
	var userBots []UserBot
	var hashes []string
	db.Conn.Table("user_bots").Select("hash").Where("user_id=?", userId).Scan(&userBots)
	for _, bot := range userBots {
		hashes = append(hashes, bot.Hash)
	}
	return hashes
}

func GetUserBotByUserId(userId string) []*BotInfoRes {
	var bots []Bot
	db.Conn.Raw("SELECT bots.client_id, bots.full_name, bots.identity_number, bots.avatar_url FROM user_bots LEFT JOIN bots ON user_bots.client_id=bots.client_id WHERE user_id=$1", userId).Scan(&bots)
	respBots := make([]*BotInfoRes, 0)
	for _, bot := range bots {
		respBots = append(respBots, &BotInfoRes{
			FullName:       bot.FullName,
			IdentityNumber: bot.IdentityNumber,
			AvatarURL:      bot.AvatarURL,
			ClientId:       bot.ClientId,
		})
	}
	return respBots
}

func CheckUserHasBot(userId, clientId string) *Bot {
	var userBot UserBot
	db.Conn.First(&userBot, "user_id=? AND client_id=?", userId, clientId)
	if userBot.ClientId == "" {
		return nil
	}
	var bot Bot
	db.Conn.First(&bot, "client_id=? AND hash=?", userBot.ClientId, userBot.Hash)
	if bot.Hash == "" {
		return nil
	}
	return &bot
}

func DeleteUserBotItem(userId, clientId string) {
	db.Conn.Delete(UserBot{}, "user_id=? AND client_id=?", userId, clientId)
}

func AddOrUpdateBotItem(clientId, sessionId, privateKey, fullName, identityNumber, avatarURL string) {
	hash := Sha256Hash(clientId, sessionId, privateKey)
	var bot = Bot{
		ClientId:       clientId,
		SessionId:      sessionId,
		PrivateKey:     privateKey,
		FullName:       fullName,
		IdentityNumber: identityNumber,
		AvatarURL:      avatarURL,
		Hash:           hash,
	}
	updateStr := fmt.Sprintf(
		"ON CONFLICT(%s) DO UPDATE SET session_id='%s', private_key='%s', hash='%s', full_name='%s', identity_number='%s', avatar_url='%s'",
		"client_id", sessionId, privateKey, hash, fullName, identityNumber, avatarURL)
	db.Conn.Set("gorm:insert_option", updateStr).Create(&bot)
}

func GetAllBot() []UserBot {
	var allBot []UserBot
	db.Conn.Table("bots").Select("client_id, session_id, private_key, hash").Scan(&allBot)
	return allBot
}

func GetBotById(clientId string) Bot {
	var bot Bot
	db.Conn.First(&bot, "client_id=?", clientId)
	return bot
}

func GetBotByHash(hashList []string) []UserBase {
	var bot []UserBase
	db.Conn.Table("bots").Select("full_name,identity_number,avatar_url").Where("hash IN (?)", hashList).Scan(&bot)
	return bot
}

func DeleteBotItem(clientId string) {
	db.Conn.Delete(Bot{}, "client_id=?", clientId)
}

func Sha256Hash(clientId, sessionId, privateKey string) string {
	str := strings.Join([]string{clientId, sessionId, privateKey}, ",")
	h := sha256.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	s := hex.EncodeToString(sum)
	return string(s)
}
