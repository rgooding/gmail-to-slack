package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rgooding/gmail-to-slack/config"
)

func Send(channel, sender, subject, body string) error {
	bodyParts := chunkMessage(body)

	n := len(bodyParts)
	if n > 1 {
		bodyParts[0] = subject + fmt.Sprintf(" _(1/%d)_\n", n) + "```" + bodyParts[0] + "```"
		for i := 1; i < n; i++ {
			bodyParts[i] = subject + fmt.Sprintf(" _(%d/%d)_\n", i+1, n) + "```" + bodyParts[i] + "```"
		}
	} else {
		bodyParts[0] = subject + "\n```" + bodyParts[0] + "```"
	}

	for _, m := range bodyParts {
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
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading HTTP response body: %s", err.Error())
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("received HTTP response code %d: %s", res.StatusCode, content)
	}
	return nil
}

func chunkMessage(s string) []string {
	const chunkSize = 2400
	const maxJitter = 100

	l := len(s)
	if l <= chunkSize {
		return []string{s}
	}
	var chunks []string
	start := 0
	for start < l {
		end := start + chunkSize
		if end >= l-1 {
			end = l
		} else {
			// work backwards to find the last newline
			newEnd := end
			found := false
			for !found && end-newEnd < maxJitter {
				if s[newEnd-1] == '\n' {
					found = true
				} else {
					newEnd--
				}
			}
			end = newEnd
		}
		chunks = append(chunks, s[start:end])
		start = end
	}
	return chunks
}
