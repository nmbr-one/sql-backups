package discord

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Color int

const (
	ColorRed    Color = 15548997
	ColorYellow Color = 16705372
	ColorGreen  Color = 5763719
)

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Embed struct {
	Title  string  `json:"title"`
	Fields []Field `json:"fields"`
	Color  Color   `json:"color"`
}

type Message struct {
	Embeds []Embed `json:"embeds"`
}

func SendMessageWithFileToWebhook(webhookURL string, message Message, fileName string) error {
	var body = &bytes.Buffer{}

	var multipartWriter = multipart.NewWriter(body)

	// JSON START
	var jsonPayload, mErr = json.Marshal(message)

	if mErr != nil {
		return mErr
	}

	var jsonWriter, _ = multipartWriter.CreateFormField("payload_json")

	jsonWriter.Write(jsonPayload)
	// JSON END

	// FILE START
	var file, openErr = os.Open(fileName)

	if openErr != nil {
		return openErr
	}

	var partWriter, _ = multipartWriter.CreateFormFile("file1", fileName)

	io.Copy(partWriter, file)

	file.Close()
	// FILE END

	multipartWriter.Close()

	var req, reqErr = http.NewRequest("POST", webhookURL, body)

	if reqErr != nil {
		return reqErr
	}

	req.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	var _, resErr = http.DefaultClient.Do(req)

	if resErr != nil {
		return resErr
	}

	return nil
}

func SendMessageToWebhook(webhookURL string, message Message) error {
	var reqBody, mErr = json.Marshal(message)

	if mErr != nil {
		return mErr
	}

	var req, reqErr = http.NewRequest("POST", webhookURL, bytes.NewReader(reqBody))

	if reqErr != nil {
		return reqErr
	}

	req.Header.Add("Content-Type", "application/json")

	var _, resErr = http.DefaultClient.Do(req)

	if resErr != nil {
		return resErr
	}

	return nil
}
