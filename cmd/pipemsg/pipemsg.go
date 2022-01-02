package main

// pipemsg tool receives a raw email message on stdin and sends it to a Slack channel
// The Slack URL is read from config.yaml in the current directory (or CONFIG_FILE env var)
// The Slack channel is passed as a single argument to the executable

import (
	"encoding/base64"
	"fmt"
	"github.com/rgooding/gmail-to-slack/slack"
	"io/ioutil"
	"mime/quotedprintable"
	"net/mail"
	"os"
	"strings"
)

func fatal(f string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, f+"\n", v...)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		fatal("Usage: %s slack-channel", os.Args[0])
	}
	channel := os.Args[1]

	// Read the message content from stdin
	msg, err := mail.ReadMessage(os.Stdin)
	if err != nil {
		fatal("Error reading message from stdin: %s", err.Error())
	}

	from := msg.Header.Get("From")
	subject := msg.Header.Get("Subject")

	// Check the encoding type and decode the body if required
	var body []byte
	switch strings.ToLower(msg.Header.Get("Content-Transfer-Encoding")) {
	case "quoted-printable":
		body, err = ioutil.ReadAll(quotedprintable.NewReader(msg.Body))
	case "base64":
		body, err = ioutil.ReadAll(base64.NewDecoder(base64.URLEncoding, msg.Body))
	default:
		body, err = ioutil.ReadAll(msg.Body)
	}
	if err != nil {
		fatal("Error reading message body: %s", err.Error())
	}

	err = slack.Send(channel, from, subject, string(body))
	if err != nil {
		fatal("Error sending Slack message: %s", err.Error())
	}
}
