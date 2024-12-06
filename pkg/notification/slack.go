package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

type SlackNotification struct {
	conf *settings.Settings
}

func NewSlackNotification(conf *settings.Settings) INotification {
	return &SlackNotification{conf: conf}
}

func (n *SlackNotification) Send(domain, currentIP string) (err error) {
	if n.conf.Notify.Slack.BotAPIToken == "" {
		return errors.New("bot api token cannot be empty")
	}

	if n.conf.Notify.Slack.Channel == "" {
		return errors.New("channel cannot be empty")
	}

	client := utils.GetHTTPClient(n.conf)
	tpl := n.conf.Notify.Slack.MsgTemplate
	if tpl == "" {
		tpl = "_Your IP address is changed to_\n\n*{{ .CurrentIP }}*\n\nDomain *{{ .Domain }}* updated"
	}

	msg := buildTemplate(currentIP, domain, tpl)

	var response *http.Response
	formData := url.Values{
		"token":   {n.conf.Notify.Slack.BotAPIToken},
		"channel": {n.conf.Notify.Slack.Channel},
		"text":    {msg},
	}
	response, err = client.PostForm("https://slack.com/api/chat.postMessage", formData)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	type ResponseParameters struct {
		MigrateToChatID int64 `json:"migrate_to_chat_id"` // optional
		RetryAfter      int   `json:"retry_after"`        // optional
	}

	type APIResponse struct {
		Ok          bool                `json:"ok"`
		Result      json.RawMessage     `json:"result"`
		ErrorCode   int                 `json:"error_code"`
		Description string              `json:"description"`
		Parameters  *ResponseParameters `json:"parameters"`
	}

	var resp APIResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		fmt.Println("error:", err)
		return errors.New("failed to parse response")
	}

	if !resp.Ok {
		return errors.New(resp.Description)
	}

	return nil
}
