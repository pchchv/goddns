package notification

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

type PushoverNotification struct {
	conf *settings.Settings
}

func NewPushoverNotification(conf *settings.Settings) INotification {
	return &PushoverNotification{conf: conf}
}

func (n *PushoverNotification) Send(domain, currentIP string) (err error) {
	if n.conf.Notify.Pushover.Token == "" {
		return errors.New("pushover api token cannot be empty")
	}

	if n.conf.Notify.Pushover.User == "" {
		return errors.New("pushover user cannot be empty")
	}

	client := utils.GetHTTPClient(n.conf)
	tpl := n.conf.Notify.Pushover.MsgTemplate
	if tpl == "" {
		tpl = "Your IP address changed to <b>{{ .CurrentIP }}</b>. The DNS record for {{ .Domain }} updated."
		n.conf.Notify.Pushover.HTML = 1
	}

	var response *http.Response
	msg := buildTemplate(currentIP, domain, tpl)

	form := url.Values{}
	form.Add("token", n.conf.Notify.Pushover.Token)
	form.Add("user", n.conf.Notify.Pushover.User)
	form.Add("message", msg)
	form.Add("html", strconv.FormatInt(int64(n.conf.Notify.Pushover.HTML), 10))
	if n.conf.Notify.Pushover.Device != "" {
		form.Add("device", n.conf.Notify.Pushover.Device)
	}

	if n.conf.Notify.Pushover.Title != "" {
		form.Add("title", n.conf.Notify.Pushover.Title)
	}

	priority := n.conf.Notify.Pushover.Priority
	if priority != 0 {
		form.Add("priority", strconv.FormatInt(int64(priority), 10))
	}

	const ReqURL = "https://api.pushover.net/1/messages.json"
	log.Printf("Pushover api request URL: %s, Form: %v", ReqURL, form)
	response, err = client.PostForm(ReqURL, form)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	type APIResponse struct {
		Status  int      `json:"status"`
		Request string   `json:"request"`
		Errors  []string `json:"errors"`
		User    string   `json:"user"`
		Token   string   `json:"token"`
	}

	var resp APIResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Println("error:", err)
		return errors.New("failed to parse pushover api response")
	}

	log.Printf("Pushover api response: %+v", resp)

	if resp.Status != 1 {
		return errors.New("pushover api call failed Status: " + string(rune(resp.Status)) + ", Errors: " + strings.Join(resp.Errors, ", "))
	}

	return nil
}
