package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendSlackAlert(webhookURL string, message string) error {
	payload := map[string]string{"text": "🚨 *Gopher-Watch Alert* 🚨\n" + message}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
