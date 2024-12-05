package ovh

import (
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

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

func (provider *DNSProvider) recordTypeToIPType(Type string) string {
	if Type == utils.IPTypeAAAA {
		return utils.IPV6
	}

	return utils.IPV4
}
