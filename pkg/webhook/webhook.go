package webhook

import (
	"net/http"

	"github.com/pchchv/goddns/internal/settings"
)

type Webhook struct {
	conf   *settings.Settings
	client *http.Client
}
