package noip

import "github.com/pchchv/goddns/internal/settings"

type DNSProvider struct {
	configuration *settings.Settings
}

// Init passes DNS settings and store it to provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}
