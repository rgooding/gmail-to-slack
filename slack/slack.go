package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rgooding/gmail-to-slack/config"
)

// Messages longer than chunkSize will be split into multiple messages of this length
const chunkSize = 2400

// Messages longer than maxMessageSize will be truncated
const maxMessageSize = 10 * chunkSize

func Send(channel, sender, subject, body string) error {
	if len(body) > maxMessageSize {
		chunks := chunkMessage(body)
		msg := subject + "\n" +
			"_Message truncated. See email for the full contents_\n" +
			"```" + chunks[0] + "```" +
			"\n_...truncated..._\n"
		return sendMsg(channel, sender, msg)
	}

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

	var sleepTime time.Duration
	n := 0
	for {
		statusCode, content, err := sendRequest(jsonPayload)
		if err != nil {
			return err
		}
		if statusCode == 429 {
			n++
			if n > 10 {
				return errors.New("too many retries")
			}
			// work out how long to sleep, default to 2 seconds
			sleepTime = 2 * time.Second
			var resp map[string]interface{}
			err := json.Unmarshal(content, &resp)
			if err != nil {
				log.Printf("Error unmarshalling response from Slack: %s", err.Error())
			} else if secs, ok := resp["retry_after"].(int); ok && secs > 0 {
				sleepTime = time.Duration(secs) * time.Second
			}
			log.Printf("Slack rate-limited, sleeping %d seconds", int(sleepTime.Seconds()))
			time.Sleep(sleepTime)
			sleepTime *= 2
		} else if statusCode < 200 || statusCode > 299 {
			return fmt.Errorf("received HTTP response code %d: %s", statusCode, content)
		} else {
			break
		}
	}
	return nil
}

func sendRequest(jsonPayload []byte) (int, []byte, error) {
	cfg := config.Load()
	res, err := http.Post(cfg.SlackUrl, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, fmt.Errorf("error reading HTTP response body: %s", err.Error())
	}
	return res.StatusCode, content, nil
}

func chunkMessage(s string) []string {
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
