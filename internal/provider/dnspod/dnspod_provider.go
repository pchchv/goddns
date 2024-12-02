package dnspod

import "github.com/pchchv/goddns/internal/settings"

type DNSProvider struct {
	configuration *settings.Settings
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}
