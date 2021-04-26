package controllers

import (
	"github.com/MixinNetwork/bot-manager/models"
	"github.com/MixinNetwork/bot-manager/session"
	"github.com/astaxie/beego/context"
)

func checkBotManager(userId, clientId string, ctx *context.Context) bool {
	if clientId == "" {
		err := session.ForbiddenError()
		session.HandleError(ctx, err)
		return false
	}
	if models.CheckUserHasBot(userId, clientId) == nil {
		err := session.ForbiddenError()
		session.HandleError(ctx, err)
		return false
	}
	return true
}
