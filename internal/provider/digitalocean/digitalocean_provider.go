package digitalocean

import "github.com/pchchv/goddns/internal/settings"

type DNSProvider struct {
	configuration *settings.Settings
	API           string
}
