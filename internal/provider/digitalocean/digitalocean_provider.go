package digitalocean

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

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

// newRequest creates a new request with auth in place and optional proxy.
func (provider *DNSProvider) newRequest(method, url string, body io.Reader) (*http.Request, *http.Client) {
	client := utils.GetHTTPClient(provider.configuration)
	if client == nil {
		log.Print("cannot create HTTP client")
	}

	req, _ := http.NewRequest(method, provider.API+url, body)
	req.Header.Set("Content-Type", "application/json")

	if provider.configuration.Email != "" && provider.configuration.Password != "" {
		req.Header.Set("X-Auth-Email", provider.configuration.Email)
		req.Header.Set("X-Auth-Key", provider.configuration.Password)
	} else if provider.configuration.LoginToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.configuration.LoginToken))
	}

	log.Printf("Created %+v request for %+v", string(method), string(url))

	return req, client
}
