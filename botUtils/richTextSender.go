package botUtils

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/slainsama/msgr_server/globals"
	"github.com/slainsama/msgr_server/models"
)

func SendTextMessage(msg models.Message) {
	config := globals.UnmarshaledConfig

	message := msg
	url := config.Bot.APIUrl + config.Bot.Token + config.Bot.Methods.SendMessage
	params := map[string]string{
		"chat_id": strconv.Itoa(message.ChatId),
		"text":    message.Data,
	}
	reqURL := buildURL(url, params)
	response, err := http.Get(reqURL)
	if err != nil {
		log.Println("Error sending GET request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing :", err)
		}
	}(response.Body)
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}
	log.Println("Response Body:", buf.String())
}

/*
func SendPhotoMessage(msg string) {
	message := handleMsg(msg)
	url := globals.UnmarshaledConfig.Bot.Token + "sendPhoto"
	params := map[string]string{
		"chat_id": message.ChatId,
		"photo":   message.Photo,
	}
	reqURL := buildURL(url, params)
	response, err := http.Get(reqURL)
	if err != nil {
		log.Println("Error sending GET request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing :", err)
		}
	}(response.Body)
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}
	log.Println("Response Body:", buf.String())
}
*/

func buildURL(baseURL string, params map[string]string) string {
	url := baseURL + "?"
	for key, value := range params {
		url += key + "=" + value + "&"
	}
	url = url[:len(url)-1]
	return url
}

/*
func handleMsg(msg string) *models.Message {
	decodedMsgBytes, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		log.Println("Error decode base64:", err)
		return nil
	}
	msg = string(decodedMsgBytes)
	var message = new(models.Message)
	err = xml.Unmarshal([]byte(msg), &message)
	if err != nil {
		log.Println("Error unmarshalling XML:", err)
		return nil
	}
	return message
}
*/
