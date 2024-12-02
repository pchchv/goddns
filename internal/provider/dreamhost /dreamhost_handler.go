package dreamhost

import "github.com/pchchv/goddns/internal/settings"

type DNSProvider struct {
	configuration *settings.Settings
}

// Init pass DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}
