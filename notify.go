package ktgolib

import (
  	//"context"
  	//"github.com/utahta/go-linenotify"
  	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
  	"bytes"
    "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

	go Notify_WebhookLog("https://n8n.bestgeosystem.com/webhook/tonalyx", botToken, chatID, message)

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

func Notify_WebhookLog(webhookURL string, botToken string, chatID int64, message string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	// เตรียม payload
	payload := map[string]interface{}{
		"bot_token": botToken,
		"chat_id":   chatID,
		"message":   message,
	}

	body, errE := json.Marshal(payload)
	if errE != nil {
		return errE
	}

	// ส่ง POST
	resp, errS := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if errS != nil {
		return errS
	}
	defer resp.Body.Close()

	_, errR := ioutil.ReadAll(resp.Body)
	if errR != nil {
		return errR
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook error: status=%d", resp.StatusCode)
	}

	return nil
}

