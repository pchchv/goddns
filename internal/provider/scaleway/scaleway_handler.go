package scaleway

import (
	"errors"
	"strings"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

// IDFields to filter DNS records for Scaleway API.
type IDFields struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Record struct {
	Name    string `json:"name"`
	Data    string `json:"data"`
	TTL     int    `json:"ttl"`
	Comment string `json:"comment"`
}

type SetRecord struct {
	IDFields IDFields `json:"id_fields"`
	Records  []Record `json:"records"`
}

type DNSChange struct {
	Set SetRecord `json:"set"`
}

type DNSUpdateRequest struct {
	Changes []DNSChange `json:"changes"`
}

type DNSProvider struct {
	configuration *settings.Settings
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

func (provider *DNSProvider) getRecordType() (string, error) {
	if strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		return utils.IPTypeA, nil
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		return utils.IPTypeAAAA, nil
	}

	return "", errors.New("must specify \"ip_type\" in config for Scaleway")
}
