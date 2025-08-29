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
	"strings"
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

// ถ้า video == "" -> ส่งข้อความ HTML ปกติ
// ถ้า video != "" -> ส่ง Video พร้อม Caption (HTML)
func Notify_ToTelegramWithHTMLAndVideo(botToken string, chatID int64, messageHTML string, video string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	// ถ้าไม่ได้แนบวิดีโอ -> ส่งข้อความปกติ
	if strings.TrimSpace(video) == "" {
		msg := tgbotapi.NewMessage(chatID, messageHTML)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.DisableWebPagePreview = true // ค่า default คือ false (ให้โชว์ preview link)
		_, err = bot.Send(msg)
		return err
	}

	// ถ้ามีวิดีโอ -> ส่งเป็น Video พร้อม caption
	var v tgbotapi.VideoConfig
	if strings.HasPrefix(video, "http://") || strings.HasPrefix(video, "https://") {
		v = tgbotapi.NewVideo(chatID, tgbotapi.FileURL(video))
	} else {
		v = tgbotapi.NewVideo(chatID, tgbotapi.FilePath(video))
	}

	v.Caption = messageHTML
	v.ParseMode = tgbotapi.ModeHTML
	v.SupportsStreaming = true

	_, err = bot.Send(v)
	return err
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

