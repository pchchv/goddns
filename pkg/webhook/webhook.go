package webhook

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
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

func (w *Webhook) Execute(domain, currentIP string) (err error) {
	if w.conf.Webhook.URL == "" {
		log.Print("Webhook URL is empty, skip sending notification")
		return nil
	}

	// set request method
	method := http.MethodGet
	if w.conf.Webhook.RequestBody != "" {
		method = http.MethodPost
	}

	// send HTTP get request
	var reqURL, reqBody string
	if method == http.MethodGet {
		reqURL, err = w.buildReqURL(domain, currentIP, w.conf.IPType)
		if err != nil {
			return err
		}
	} else {
		reqURL = w.conf.Webhook.URL
		reqBody, err = w.buildReqBody(domain, currentIP, w.conf.IPType)
		if err != nil {
			return err
		}
	}

	var req *http.Request
	if req, err = http.NewRequest(method, reqURL, strings.NewReader(reqBody)); err != nil {
		log.Fatal("Failed to create request:", err)
		return err
	}

	if method == http.MethodPost {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := w.client.Do(req)
	if err != nil {
		log.Fatal("Failed to send request:", err)
		return err
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
		return err
	}

	log.Printf("Webhook response: %s", string(content))
	return nil
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

func (w *Webhook) buildReqURL(domain, currentIP, ipType string) (string, error) {
	t := template.New("req template")
	if _, err := t.Parse(w.conf.Webhook.URL); err != nil {
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
