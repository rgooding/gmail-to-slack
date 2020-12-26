package slack

import (
	"bytes"
	"encoding/json"
	"github.com/rgooding/gmail-to-slack/config"
	"log"
	"net/http"
)

func Send(Channel, Sender, Subject, Body string) error {
	payload := map[string]interface{}{
		"channel":      "#" + Channel,
		"icon_emoji":   ":information_source:",
		"username":     Subject,
		"text":         Body,
		"unfurl_links": false,
		"unfurl_media": false,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	cfg := config.Load()
	_, err = http.Post(cfg.SlackUrl, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		log.Printf("Error sending Slack notification: %s", err.Error())
	}
	return err
}
