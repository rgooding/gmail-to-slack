package main

import (
	"github.com/rgooding/gmail-to-slack/config"
	"github.com/rgooding/gmail-to-slack/gmailclient"
	"github.com/rgooding/gmail-to-slack/slack"
	"log"
	"strings"
)

func main() {
	cfg := config.Load()

	labelNames := mapKeys(cfg.LabelChannels)
	labelIds, err := gmailclient.GetLabelIds(labelNames)
	if err != nil {
		log.Fatal(err)
	}

	// build a map of label ID => slack channel
	labelIdChannels := make(map[string]string)
	for labelName, channel := range cfg.LabelChannels {
		for id, name := range labelIds {
			if labelName == name {
				labelIdChannels[id] = channel
				break
			}
		}
	}

	log.Printf("Finding messages with label(s): %s", strings.Join(labelNames, ", "))
	messages, err := gmailclient.ListUnreadMessages(mapKeys(labelIds))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d new messages", len(messages))

	var sentIds []string
	for _, msg := range messages {
		for _, labelId := range msg.LabelIds {
			if channel, ok := labelIdChannels[labelId]; ok {
				log.Printf("Sending message %s to Slack #%s, Subject: %s", msg.Id, channel, msg.Subject)
				err = slack.Send(channel, msg.Sender, msg.Subject, msg.Body)
				if err != nil {
					log.Printf("Error sending Slack message: %s", err.Error())
				} else {
					sentIds = append(sentIds, msg.Id)
				}
			}
		}
	}

	if len(sentIds) > 0 {
		err := gmailclient.MarkAsRead(sentIds, cfg.ArchiveSentMessages)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func mapKeys(srcMap map[string]string) []string {
	keys := make([]string, 0, len(srcMap))
	for k := range srcMap {
		keys = append(keys, k)
	}
	return keys
}
