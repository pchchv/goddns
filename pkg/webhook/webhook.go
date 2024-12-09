package webhook

import (
	"bytes"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

var (
	instance *Webhook
	once     sync.Once
)

type Webhook struct {
	conf   *settings.Settings
	client *http.Client
}

func (w *Webhook) buildReqBody(domain, currentIP, ipType string) (string, error) {
	t := template.New("reqBody template")
	if _, err := t.Parse(w.conf.Webhook.RequestBody); err != nil {
		log.Fatal("Failed to parse template:", err)
		return "", err
	}

	data := struct {
		CurrentIP string
		Domain    string
		IPType    string
	}{
		currentIP,
		domain,
		ipType,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Fatal(err)
		return "", err
	}

	return tpl.String(), nil
}

func GetWebhook(conf *settings.Settings) *Webhook {
	once.Do(func() {
		instance = &Webhook{
			conf:   conf,
			client: utils.GetHTTPClient(conf),
		}
	})

	return instance
}
