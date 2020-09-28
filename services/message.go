package services

import (
	"github.com/liuzemei/bot-manager/externals"
	"github.com/liuzemei/bot-manager/models"
	"log"
	"time"
)

var HashMap = map[string]bool{}

func ConnectedBot() {
	for {
		checkConnect()
		time.Sleep(time.Minute * 5)
	}
}

func checkConnect() {
	bots := models.GetAllBot()
	for _, bot := range bots {
		if !HashMap[bot.Hash] {
			go connectBot(bot)
		}
	}
}

func connectBot(botInfo models.UserBot) {
	HashMap[botInfo.Hash] = true
	err := externals.StartWebSockets(botInfo.ClientId, botInfo.SessionId, botInfo.PrivateKey, botInfo.Hash)
	if err != nil {
		delete(HashMap, botInfo.Hash)
		log.Println("开启WebSocket失败！", err)
		//models.DeleteBotItem(botInfo.ClientId)
	}
}
