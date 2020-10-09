package controllers

import (
	"github.com/astaxie/beego/context"
	"github.com/liuzemei/bot-manager/models"
	"github.com/liuzemei/bot-manager/session"
)

func checkBotManager(userId, clientId string, ctx *context.Context) bool {
	if clientId == "" {
		err := session.ForbiddenError()
		session.HandleError(ctx, err)
		return false
	}
	if models.CheckUserHasBot(userId, clientId) {
		err := session.ForbiddenError()
		session.HandleError(ctx, err)
		return false
	}
	return true
}
