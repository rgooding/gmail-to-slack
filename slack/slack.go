package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rgooding/gmail-to-slack/config"
	"io/ioutil"
	"net/http"
)

func Send(channel, sender, subject, body string) error {
	fullMsg := subject + "\n```" + body + "```"
	messages := chunkString(fullMsg, 3900)
	l := len(messages)
	if l > 1 {
		messages[0] += "```"
		for i := 1; i < l-1; i++ {
			messages[i] = "```" + messages[i] + "```"
		}
		messages[l-1] = "```" + messages[l-1]
	}
	for _, m := range messages {
		err := sendMsg(channel, sender, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendMsg(channel, sender, body string) error {
	payload := map[string]interface{}{
		"channel":      "#" + channel,
		"icon_emoji":   ":information_source:",
		"username":     sender,
		"text":         body,
		"unfurl_links": false,
		"unfurl_media": false,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	cfg := config.Load()
	res, err := http.Post(cfg.SlackUrl, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading HTTP response body: %s", err.Error())
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("received HTTP response code %d: %s", res.StatusCode, content)
	}
	return nil
}

func chunkString(s string, chunkSize int) []string {
	l := len(s)
	if l <= chunkSize {
		return []string{s}
	}

	var chunks []string
	for i := 0; i < l; i += chunkSize {
		end := i + chunkSize
		if end > l-1 {
			end = l
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
