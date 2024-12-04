package hetzner

import (
	"net/http"

	"github.com/pchchv/goddns/internal/settings"
)

type Record struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    int64  `json:"ttl"`
	ZoneID string `json:"zone_id"`
}

type DNSProvider struct {
	configuration *settings.Settings
	client        *http.Client
}
