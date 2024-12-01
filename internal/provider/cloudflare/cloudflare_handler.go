package cloudflare

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

// URL is the endpoint for the Cloudflare API.
const URL = "https://api.cloudflare.com/client/v4"

// Zone object with id and name.
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ZoneResponse is a wrapper for Zones.
type ZoneResponse struct {
	Zones   []Zone `json:"result"`
	Success bool   `json:"success"`
}

// DNSRecord for Cloudflare API.
type DNSRecord struct {
	ID      string `json:"id"`
	IP      string `json:"content"`
	Name    string `json:"name"`
	Proxied bool   `json:"proxied"`
	Type    string `json:"type"`
	ZoneID  string `json:"zone_id"`
	TTL     int32  `json:"ttl"`
}

// SetIP updates DNSRecord.IP.
func (r *DNSRecord) SetIP(ip string) {
	r.IP = ip
}

type DNSRecordUpdateResponse struct {
	Record  DNSRecord `json:"result"`
	Success bool      `json:"success"`
}

type DNSRecordResponse struct {
	Records []DNSRecord `json:"result"`
	Success bool        `json:"success"`
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
		log.Println("cannot create HTTP client")
	}

	req, _ := http.NewRequest(method, provider.API+url, body)
	req.Header.Set("Content-Type", "application/json")

	if provider.configuration.Email != "" && provider.configuration.Password != "" {
		req.Header.Set("X-Auth-Email", provider.configuration.Email)
		req.Header.Set("X-Auth-Key", provider.configuration.Password)
	} else if provider.configuration.LoginToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.configuration.LoginToken))
	}

	return req, client
}

// getZone find the correct zone via domain name.
func (provider *DNSProvider) getZone(domain string) string {
	var z ZoneResponse

	req, client := provider.newRequest("GET", fmt.Sprintf("/zones?name=%s", domain), nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request error:", err)
		return ""
	}

	body, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &z); err != nil {
		log.Fatalf("Decoder error: %+v", err)
		log.Printf("Response body: %+v", string(body))
		return ""
	} else if !z.Success {
		log.Printf("Response failed: %+v", string(body))
		return ""
	}

	for _, zone := range z.Zones {
		if zone.Name == domain {
			return zone.ID
		}
	}

	return ""
}

func (provider *DNSProvider) getCurrentDomain(domainName string) *settings.Domain {
	for _, domain := range provider.configuration.Domains {
		domain := domain
		if domain.DomainName == domainName {
			return &domain
		}
	}

	return nil
}

// Get all DNS A records for a zone.
func (provider *DNSProvider) getDNSRecords(zoneID string) []DNSRecord {
	var empty []DNSRecord
	var recordType string
	var r DNSRecordResponse
	if provider.configuration.IPType == "" || strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		recordType = utils.IPTypeA
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		recordType = utils.IPTypeAAAA
	}

	log.Printf("Querying records with type: %s", recordType)
	req, client := provider.newRequest("GET", fmt.Sprintf("/zones/"+zoneID+"/dns_records?type=%s&page=1&per_page=500", recordType), nil)
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
	} else if !r.Success {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Response failed: %+v", string(body))
		return empty

	}

	return r.Records
}
