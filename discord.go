package main

import (
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

type WebhookPayload struct {
	Content string         `json:"content,omitempty"`
	Embeds  []WebhookEmbed `json:"embeds"`
}

type WebhookEmbed struct {
	Title       string              `json:"title"`
	Type        string              `json:"type"`
	Description string              `json:"description"`
	Color       int                 `json:"color,omitempty"` // https://gist.github.com/thomasbnt/b6f455e2c7d743b796917fa3c205f812?permalink_comment_id=3656937#gistcomment-3656937
	Author      WebhookEmbedAuthor  `json:"author,omitempty"`
	Fields      []WebhookEmbedField `json:"fields,omitempty"`
	URL         string              `json:"url,omitempty"`
}

type WebhookEmbedAuthor struct {
	Name string `json:"name"`
}

type WebhookEmbedField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func invokeWebhook(payload WebhookPayload) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(os.Getenv("DISCORD_WEBHOOK_URL") + "?wait=true")

	if err != nil {
		log.Fatalln(err)
	} else if resp.IsError() {
		log.Fatalln(resp)
	}
}
