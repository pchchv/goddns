package webhook

import (
	"net/http"
	"sync"

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

func GetWebhook(conf *settings.Settings) *Webhook {
	once.Do(func() {
		instance = &Webhook{
			conf:   conf,
			client: utils.GetHTTPClient(conf),
		}
	})

	return instance
}
