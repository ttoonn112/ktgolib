package ktgolib

import (
  //"bytes"
  //"context"
  //"github.com/utahta/go-linenotify"
  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
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
}*/

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

func Notify_ToTelegramWithHTML(botToken string, chatID int64, message string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}
	
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

