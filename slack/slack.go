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
	payload := map[string]interface{}{
		"channel":    "#" + channel,
		"icon_emoji": ":information_source:",
		"username":   sender,
		"text":       subject,
		"attachments": []map[string]interface{}{
			{
				"id":   1,
				"text": "```" + body + "```",
			},
		},
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
