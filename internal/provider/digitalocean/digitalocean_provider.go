package digitalocean

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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

func (provider *DNSProvider) getRecordType() string {
	var recordType string = utils.IPTypeA
	if provider.configuration.IPType == "" || strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		recordType = utils.IPTypeA
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		recordType = utils.IPTypeAAAA
	}

	return recordType
}

// getDNSRecords gets all DNS A(AAA) records for a zone.
func (provider *DNSProvider) getDNSRecords(domainName string) []DNSRecord {
	var empty []DNSRecord
	var r DomainRecordsResponse
	recordType := provider.getRecordType()

	log.Printf("Querying records with type: %s", recordType)
	req, client := provider.newRequest("GET", fmt.Sprintf("/domains/"+domainName+"/records?type=%s&page=1&per_page=200", recordType), nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request error:", err)
		return empty
	}

	body, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &r); err != nil {
		log.Printf("Decoder error: %+v", err)
		log.Printf("Response body: %+v", string(body))
		return empty
	}

	return r.Records
}
