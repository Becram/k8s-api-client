package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type SlackRequestBody struct {
	Text     string `json:"text"`
	Channel  string `json:"channel"`
	Username string `json:"username"`
}

// func main() {
// 	SendSlackNotification("test-svc", "Test Message from golangcode.com")
// }

func appInit() {

	if os.Getenv("SLACK_WEBHOOK_URL") == "" {
		fmt.Fprintln(os.Stderr, "Please provide 'SLACK_WEBHOOK_URL' through environment")
		os.Exit(1)
	}

	if os.Getenv("SLACK_CHANNEL") == "" {
		fmt.Fprintln(os.Stderr, "Please provide 'SLACK_CHANNEL' through environment")
		os.Exit(1)
	}

}

func SendSlackNotification(svc string, msg string) error {

	appInit()

	webhookUrl := os.Getenv("SLACK_WEBHOOK_URL")
	username := os.Getenv("SLACK_USERNAME")
	channel := os.Getenv("SLACK_CHANNEL")
	message := svc + ":  " + msg

	slackBody, _ := json.Marshal(SlackRequestBody{Text: message, Channel: channel, Username: username})
	fmt.Printf("slack %s", slackBody)
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("non-ok response returned from slack")
	}
	return nil
}
