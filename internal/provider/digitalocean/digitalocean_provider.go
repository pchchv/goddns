package digitalocean

import "github.com/pchchv/goddns/internal/settings"

// URL is the endpoint for the DigitalOcean API.
const URL = "https://api.digitalocean.com/v2"

// DNSRecord for DigitalOcean API.
type DNSRecord struct {
	ID   int32  `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	IP   string `json:"data"`
	TTL  int32  `json:"ttl"`
}

// SetIP updates DNSRecord.IP.
func (r *DNSRecord) SetIP(ip string) {
	r.IP = ip
}

type DomainRecordsResponse struct {
	Records []DNSRecord `json:"domain_records"`
}

type DNSProvider struct {
	configuration *settings.Settings
	API           string
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
	provider.API = URL
}
