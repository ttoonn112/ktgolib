package ktgolib

import (
  //"bytes"
	"context"
  "github.com/utahta/go-linenotify"
  "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Notify_ToLine(token string, msg string) (string, string) {
	c := linenotify.NewClient()
	response, err := c.Notify(context.Background(), token, msg, "", "", nil)
	if err != nil {
		return "", err.Error()
	}else if response.Status != 200 {
		return I_S(response.RateLimit.Remaining)+"/"+I_S(response.RateLimit.Limit), I_S(response.Status)+" - "+response.Message
	}else{
		return I_S(response.RateLimit.Remaining)+"/"+I_S(response.RateLimit.Limit), ""
	}

	//c.Notify(context.Background(), token, "hello world", "", "", nil)
	//c.Notify(context.Background(), token, "hello world", "http://localhost/thumb.jpg", "http://localhost/full.jpg", nil)
	//c.Notify(context.Background(), token, "hello world", "", "", bytes.NewReader([]byte("image bytes")))
}

func Notify_ToTelegram(botToken string, chatID int64, message string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(chatID, message)
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func Notify_TelegramGetUpdates(botToken string) map[string]interface{} {
	url := "https://api.telegram.org/bot"+botToken+"/getUpdates"
	result := Http_Get(url)
	return result
}