package digitalocean

import "github.com/pchchv/goddns/internal/settings"

// URL is the endpoint for the DigitalOcean API.
const URL = "https://api.digitalocean.com/v2"

type DNSProvider struct {
	configuration *settings.Settings
	API           string
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
	provider.API = URL
}
