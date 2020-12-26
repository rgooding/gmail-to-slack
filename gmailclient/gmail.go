package gmailclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/rgooding/gmail-to-slack/config"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"strings"
)

const scope = "https://www.googleapis.com/auth/gmail.modify"
const userId = "me"

type ParsedMessage struct {
	Id       string
	Sender   string
	Subject  string
	Body     string
	LabelIds []string
}

func (m ParsedMessage) HasLabel(labelId string) bool {
	for _, lbl := range m.LabelIds {
		if lbl == labelId {
			return true
		}
	}
	return false
}

var sharedSvc *gmail.Service

func getSvc() (*gmail.Service, error) {
	if sharedSvc == nil {
		cfg := config.Load()
		ctx := context.TODO()
		ts := NewTokenSource(cfg.SecretFile, cfg.TokenFile, scope)
		s, err := gmail.NewService(ctx, option.WithTokenSource(ts))
		if err != nil {
			return nil, err
		}
		sharedSvc = s
	}
	return sharedSvc, nil
}

func ListUnreadMessages(labelIds []string) ([]ParsedMessage, error) {
	svc, err := getSvc()
	if err != nil {
		return nil, err
	}

	messageIds := make(map[string]bool)
	for _, label := range labelIds {
		msgRes, err := svc.Users.Messages.List(userId).LabelIds(label, "UNREAD").Do()
		if err != nil {
			return nil, err
		}
		for _, msg := range msgRes.Messages {
			messageIds[msg.Id] = true
		}
	}

	var parsedMessages []ParsedMessage
	for msgId := range messageIds {
		msg, err := svc.Users.Messages.Get(userId, msgId).Do()
		if err != nil {
			return nil, err
		}
		body, err := base64.StdEncoding.DecodeString(msg.Payload.Body.Data)
		if err != nil {
			return nil, fmt.Errorf("error decoding body of message %s: %s", msgId, err.Error())
		}

		subject := ""
		sender := ""
		for _, h := range msg.Payload.Headers {
			switch strings.ToLower(h.Name) {
			case "subject":
				subject = h.Value
			case "from":
				sender = h.Value
			}
		}
		parsedMessages = append(parsedMessages, ParsedMessage{
			Id:       msgId,
			Sender:   sender,
			Subject:  subject,
			Body:     string(body),
			LabelIds: msg.LabelIds,
		})
	}
	return parsedMessages, nil
}

// Make a list of label IDs from their names
func GetLabelIds(labelNames []string) (map[string]string, error) {
	svc, err := getSvc()
	if err != nil {
		return nil, err
	}
	res, err := svc.Users.Labels.List(userId).Do()
	if err != nil {
		return nil, err
	}

	labelIds := make(map[string]string)
	for _, label := range res.Labels {
		for _, lblName := range labelNames {
			if label.Name == lblName {
				labelIds[label.Id] = label.Name
			}
		}
	}
	return labelIds, nil
}

func MarkAsRead(messageIds []string, archive bool) error {
	if len(messageIds) > 0 {
		log.Printf("Marking %d messages as read...", len(messageIds))
		svc, err := getSvc()
		if err != nil {
			return err
		}

		removeLabels := []string{"UNREAD"}
		if archive {
			removeLabels = append(removeLabels, "INBOX")
		}

		return svc.Users.Messages.BatchModify(userId, &gmail.BatchModifyMessagesRequest{
			Ids:            messageIds,
			RemoveLabelIds: removeLabels,
		}).Do()
	}
	return nil
}
