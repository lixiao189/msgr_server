package controller

import (
	"errors"

	"github.com/slainsama/msgr_server/bot/botMethod"
	botGlobals "github.com/slainsama/msgr_server/bot/globals"
	"github.com/slainsama/msgr_server/bot/handler"
	"github.com/slainsama/msgr_server/globals"
	"github.com/slainsama/msgr_server/models"
	"github.com/slainsama/msgr_server/utils"
	"gorm.io/gorm"
)

func init() {
	botGlobals.Dispatcher.AddHandler(handler.NewCommandHandler("/start", startController))
}

// startController "/start"
func startController(u *models.TelegramUpdate) {
	userInfo := u.Message.From
	var user models.User
	var message models.Message
	message.ChatId = u.Message.Chat.ID
	result := globals.DB.Where(models.User{ID: userInfo.ID}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果记录不存在，则创建新记录
			newUser := models.User{
				ID:           userInfo.ID,
				IsBot:        userInfo.IsBot,
				FirstName:    userInfo.FirstName,
				LastName:     userInfo.LastName,
				Username:     userInfo.Username,
				LanguageCode: userInfo.LanguageCode,
				IsAdmin:      false,
			}
			globals.DB.Create(&newUser)
			message.Data = utils.EscapeChar("welcome.")
			botMethod.SendTextMessage(message.ChatId, message.Data)
		}
	} else {
		message.Data = utils.EscapeChar("user already exist.")
		botMethod.SendTextMessage(message.ChatId, message.Data)
	}
}
