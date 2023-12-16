package auth

import (
	"bytes"
	"encoding/json"
	"errors"
  "io/ioutil"
	"net/http"
)

func Line_GetUserProfile(accessToken string) (map[string]interface{}, error) {
	// Line Profile API Endpoint
	apiEndpoint := "https://api.line.me/v2/profile"

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	// Set the request headers
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Create the HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error: Unexpected status code " + resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var profile map[string]interface{}
	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}

	// Check if the response contains the user ID
	userId, ok := profile["userId"].(string)
	if !ok || userId == "" {
		return nil, errors.New("Error: User ID not found in the response")
	}

	return profile, nil
}

func SendMessageToLineUser(accessToken, userId, message string) error {
	// Line Messaging API Endpoint
	apiEndpoint := "https://api.line.me/v2/bot/message/push"

	// Create the message payload
	payload := map[string]interface{}{
		"to": userId,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": message,
			},
		},
	}

	// Convert the payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Create the HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return errors.New("Error: Unexpected status code " + resp.Status)
	}

	return nil
}
