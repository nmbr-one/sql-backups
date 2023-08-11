package discord

import (
	"fmt"
)

type Logger struct {
	WebhookURL string
}

func (this Logger) LogFile(embed Embed, fileName string) {
	var message Message

	message.Embeds = append(message.Embeds, embed)

	var webhookError = SendMessageWithFileToWebhook(this.WebhookURL, message, fileName)

	if webhookError != nil {
		fmt.Printf("Sending webhook failed: %#v\n", webhookError)
	}
}

func (this Logger) Log(embed Embed) {
	var message Message

	message.Embeds = append(message.Embeds, embed)

	var webhookError = SendMessageToWebhook(this.WebhookURL, message)

	if webhookError != nil {
		fmt.Printf("Sending webhook failed: %#v\n", webhookError)
	}
}
