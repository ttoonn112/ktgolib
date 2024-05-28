package notify

import (
  "github.com/utahta/go-linenotify"
)

func Notify_ToLine(line_token string, msg string) (string, string) {
	linenotify := linenotify.New()
	response, err := linenotify.Notify(line_token, msg, "", "", nil)
	if err != nil {
		return "", err.Error()
	}else if response.Status != 200 {
		return I_S(response.RateLimit.Remaining)+"/"+I_S(response.RateLimit.Limit), I_S(response.Status)+" - "+response.Message
	}else{
		return I_S(response.RateLimit.Remaining)+"/"+I_S(response.RateLimit.Limit), ""
	}
}
