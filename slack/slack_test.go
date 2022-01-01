package slack

import (
	"github.com/rgooding/gmail-to-slack/config"
	"strconv"
	"strings"
	"testing"
)

func TestSlack(t *testing.T) {
	config.ConfFile = "../config.yaml"

	//body:="Short message body"

	// Long message
	body := ""
	for i := 0; i < 100; i++ {
		body += strconv.Itoa(i) + " " + strings.Repeat("a", 80) + "\n"
	}

	err := Send("test-channel", "Test Message", "Test Subject", body)
	if err != nil {
		t.Fatal(err)
	}
}
