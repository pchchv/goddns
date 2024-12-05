package ovh

import "github.com/pchchv/goddns/internal/settings"

type Record struct {
	Zone      string `json:"zone"`
	TTL       int    `json:"ttl"`
	Value     string `json:"target"`
	SubDomain string `json:"subDomain"`
	Type      string `json:"fieldType"`
	ID        int    `json:"id"`
}

type DNSProvider struct {
	configuration *settings.Settings
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}
