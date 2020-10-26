package controllers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/liuzemei/bot-manager/externals"
	"github.com/liuzemei/bot-manager/models"
	"github.com/liuzemei/bot-manager/utils"
	"log"
	"strings"
)

// 客户端消息入口
func readMessage(messageCome <-chan models.MessengerChannel, hash string) {
	for msg := range messageCome {
		_msg := &msg.Message
		if _msg.Source == "ACKNOWLEDGE_MESSAGE_RECEIPT" {
			models.UpdateMessage(_msg.MessageId, _msg.Status)
			for _, chanel := range models.HashManagerMap[hash] {
				if chanel != nil {
					chanel <- models.RespMessage{
						MessageId: _msg.MessageId,
						Status:    _msg.Status,
						Source:    _msg.Source,
					}
				}
			}
			continue
		}

		adminIds := models.GetAdminIdsByBotId(msg.ClientID)
		isManager := checkManager(_msg.UserId, adminIds)
		if !isManager && _msg.Category == "PLAIN_TEXT" {
			if replayed := handleAutoReplayMessage(msg, hash, adminIds); replayed {
				continue
			}
		}

		//  2. 转发给已登录的后端管理的管理员
		forwardLoginDashboardAdmins(msg, hash, isManager, adminIds)

		// 转发给手机端
		if isManager {
			// 管理员的消息
			if _msg.QuoteMessageId != "" {
				//	1. 转发给其他管理员
				//  2. 转发给被 Quote 的用户
				if hasError := handleAdminClientMessageWithQuote(msg, adminIds); !hasError {
					continue
				}
			} else {
				// 1. 判断是否是 APP_CARD
				if hasError := handleAdminClientMessageWithNoQuote(msg, adminIds); !hasError {
					continue
				}
			}
		} else {
			// 其他人员的消息
			//	1. 直接转发给管理员的 Messenger
			handleNonAdminClientMessage(msg, adminIds)
		}
	}
}

func checkManager(userId string, UserIds []string) bool {
	for _, id := range UserIds {
		if id == userId {
			return true
		}
	}
	return false
}

// 处理自动回复消息
func handleAutoReplayMessage(msg models.MessengerChannel, hash string, adminIds []string) bool {
	_msg := &msg.Message
	decodeBytes, _ := base64.StdEncoding.DecodeString(_msg.Data)
	decodeString := strings.ToLower(string(decodeBytes))
	replayData, replayCategory := models.GetAutoReplayMessageByKey(msg.ClientID, decodeString)
	if replayData != "" {
		forwardDashboardMessage(&forwardMessagePropsType{
			Category:         _msg.Category,
			CreatedAt:        utils.FormatTime(_msg.CreatedAt),
			MessageId:        _msg.MessageId,
			Source:           _msg.Source,
			UserId:           _msg.UserId,
			QuoteMessageId:   _msg.QuoteMessageId,
			ConversationId:   _msg.ConversationId,
			Data:             _msg.Data,
			RepresentativeId: _msg.RepresentativeId,
			Status:           _msg.Status,
			UpdatedAt:        utils.FormatTime(_msg.UpdatedAt),
		}, msg.ClientID, msg.SessionID, msg.PrivateKey, hash, false, adminIds)
		var sendMessages []*bot.MessageRequest
		base64Data := base64.StdEncoding.EncodeToString([]byte(replayData))
		messageId := bot.UuidNewV4().String()
		sendMessages = append(sendMessages, &bot.MessageRequest{
			MessageId:      messageId,
			Category:       replayCategory,
			Data:           base64Data,
			ConversationId: _msg.ConversationId,
			RecipientId:    _msg.UserId,
		})
		err := externals.SendBatchMessage(sendMessages, msg.ClientID, msg.SessionID, msg.PrivateKey)
		if err != nil {
			log.Println("转发自动回复的消息出问题了。", err)
			return true
		}
		forwardDashboardMessage(&forwardMessagePropsType{
			Category:       replayCategory,
			CreatedAt:      utils.FormatTime(_msg.CreatedAt),
			MessageId:      messageId,
			UserId:         _msg.UserId,
			ConversationId: _msg.ConversationId,
			Data:           base64Data,
			UpdatedAt:      utils.FormatTime(_msg.UpdatedAt),
		}, msg.ClientID, msg.SessionID, msg.PrivateKey, hash, false, adminIds)
		return true
	}
	return false
}

// 转发给管理后台
func forwardLoginDashboardAdmins(msg models.MessengerChannel, hash string, isManager bool, adminIds []string) {
	_msg := &msg.Message
	forwardDashboardMessage(&forwardMessagePropsType{
		Category:         _msg.Category,
		CreatedAt:        utils.FormatTime(_msg.CreatedAt),
		MessageId:        _msg.MessageId,
		Source:           _msg.Source,
		UserId:           _msg.UserId,
		QuoteMessageId:   _msg.QuoteMessageId,
		ConversationId:   _msg.ConversationId,
		Data:             _msg.Data,
		RepresentativeId: _msg.RepresentativeId,
		Status:           _msg.Status,
		UpdatedAt:        utils.FormatTime(_msg.UpdatedAt),
	}, msg.ClientID, msg.SessionID, msg.PrivateKey, hash, isManager, adminIds)
}

// 处理含有quote的消息客户端管理员消息
func handleAdminClientMessageWithQuote(msg models.MessengerChannel, adminIds []string) bool {
	_msg := &msg.Message
	originMessage := models.GetOriginMessageById(msg.ClientID, _msg.QuoteMessageId)
	if originMessage == nil {
		log.Println("没找到这条消息originMessage，", _msg.QuoteMessageId)
		return false
	}
	forwardMessages := models.GetForwardMessagesByOrigin(msg.ClientID, originMessage.OriginMessageId)
	if forwardMessages == nil {
		log.Println("没找到这条消息forwardMessage", originMessage.OriginMessageId)
		return false
	}
	// 转发给其他管理员，包含quote自己
	var sendMessages []*bot.MessageRequest
	for _, userId := range adminIds {
		if userId == _msg.UserId {
			continue
		}
		conversationId := bot.UniqueConversationId(msg.ClientID, userId)
		sendMessages = append(sendMessages, &bot.MessageRequest{
			MessageId:        bot.UuidNewV4().String(),
			Category:         _msg.Category,
			Data:             _msg.Data,
			QuoteMessageId:   forwardMessages[userId].MessageId,
			ConversationId:   conversationId,
			RecipientId:      userId,
			RepresentativeId: _msg.UserId,
		})
	}
	// 转发给被 Quote 的那个用户
	sendMessages = append(sendMessages, &bot.MessageRequest{
		ConversationId: originMessage.ConversationId,
		RecipientId:    originMessage.RecipientId,
		MessageId:      bot.UuidNewV4().String(),
		Category:       _msg.Category,
		Data:           _msg.Data,
		QuoteMessageId: originMessage.OriginMessageId,
	})
	err := externals.SendBatchMessage(sendMessages, msg.ClientID, msg.SessionID, msg.PrivateKey)
	if err != nil {
		log.Println("转发管理员消息出问题了。", err)
		return false
	}
	for _, message := range sendMessages {
		models.AddForwardMessage(models.ForwardMessage{
			ClientId:        msg.ClientID,
			MessageId:       message.MessageId,
			AdminId:         _msg.UserId,
			RecipientId:     message.RecipientId,
			OriginMessageId: originMessage.MessageId,
			ConversationId:  message.ConversationId,
			AdminMessageId:  _msg.MessageId,
			CreatedAt:       utils.FormatTime(_msg.CreatedAt),
		})
	}
	return true
}

// 处理管理员客户端发送的消息（无quote_id）
func handleAdminClientMessageWithNoQuote(msg models.MessengerChannel, adminIds []string) bool {
	_msg := &msg.Message
	if _msg.Category == "APP_CARD" {
		_, err := sendBroadcast(&models.Bot{
			ClientId:   msg.ClientID,
			SessionId:  msg.SessionID,
			PrivateKey: msg.PrivateKey,
		}, _msg.Category, _msg.Data, true)
		if err != nil {
			log.Println("client_message handleAdminClientMessageWithNoQuote", err)
			return false
		}
	}
	return true
}

func sendBroadcast(clientBot *models.Bot, category, data string, isBase64Data bool) ([]*bot.MessageRequest, error) {
	// 1. 拿到所有的用户
	userList := models.GetBotUserListById(clientBot.ClientId)
	// 3. 构建所有的消息
	var sendMessages []*bot.MessageRequest
	base64Data := ""
	if isBase64Data {
		base64Data = data
	} else {
		if category == "PLAIN_TEXT" {
			base64Data = base64.StdEncoding.EncodeToString([]byte(data))
		} else if category == "PLAIN_IMAGE" {
			_msgData, err := json.Marshal(data)
			base64Data = base64.StdEncoding.EncodeToString(_msgData)
			if err != nil {
				return nil, err
			}
		}
	}
	for _, user := range userList {
		sendMessages = append(sendMessages, &bot.MessageRequest{
			ConversationId: bot.UniqueConversationId(clientBot.ClientId, user),
			RecipientId:    user,
			MessageId:      bot.UuidNewV4().String(),
			Category:       category,
			Data:           base64Data,
		})
	}
	err := externals.SendBatchMessage(sendMessages, clientBot.ClientId, clientBot.SessionId, clientBot.PrivateKey)
	if err != nil {
		return nil, err
	}
	return sendMessages, nil
}

// 处理非管理员客户端消息
func handleNonAdminClientMessage(msg models.MessengerChannel, adminIds []string) {
	_msg := &msg.Message
	if _msg.Category == "MESSAGE_RECALL" {
		return
	}
	var sendMessages []*bot.MessageRequest
	for _, userId := range adminIds {
		conversationId := bot.UniqueConversationId(msg.ClientID, userId)
		messageId := bot.UuidNewV4().String()
		sendMessages = append(sendMessages, &bot.MessageRequest{
			ConversationId:   conversationId,
			RecipientId:      userId,
			RepresentativeId: _msg.UserId,
			MessageId:        messageId,
			Category:         _msg.Category,
			Data:             _msg.Data,
		})
		models.AddForwardMessage(models.ForwardMessage{
			ClientId:        msg.ClientID,
			MessageId:       messageId,
			OriginMessageId: _msg.MessageId,
			RecipientId:     _msg.UserId,
			ConversationId:  _msg.ConversationId,
			AdminId:         userId,
			CreatedAt:       utils.FormatTime(_msg.CreatedAt),
		})
	}
	if sendMessages != nil {
		err := externals.SendBatchMessage(sendMessages, msg.ClientID, msg.SessionID, msg.PrivateKey)
		if err != nil {
			log.Printf("转发普通消息出问题了。 %#v", err)
			log.Printf("%#v", sendMessages[0])
		}
	} else {
		log.Printf("没找到管理员？？%s,%#v", msg.ClientID, msg.Message)
	}
}
