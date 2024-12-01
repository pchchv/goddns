package cloudflare

import (
	"bytes"
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

// getDNSRecords gets all DNS A records for a zone.
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

func (provider *DNSProvider) createRecord(zoneID, domain, subDomain, ip string) error {
	var recordType string
	if provider.configuration.IPType == "" || strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		recordType = utils.IPTypeA
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		recordType = utils.IPTypeAAAA
	}

	newRecord := DNSRecord{
		Type: recordType,
		IP:   ip,
		TTL:  1,
	}

	if provider.configuration.Proxied {
		newRecord.Proxied = true
	}

	if subDomain == utils.RootDomain {
		newRecord.Name = utils.RootDomain
	} else {
		newRecord.Name = fmt.Sprintf("%s.%s", subDomain, domain)
	}

	content, err := json.Marshal(newRecord)
	if err != nil {
		log.Fatalf("Encoder error: %+v", err)
		return err
	}

	req, client := provider.newRequest("POST", fmt.Sprintf("/zones/%s/dns_records", zoneID), bytes.NewBuffer(content))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request error:", err)
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read request body: %+v", err)
		return err
	}

	var r DNSRecordUpdateResponse
	if err = json.Unmarshal(body, &r); err != nil {
		log.Fatalf("Decoder error: %+v", err)
		return err
	} else if !r.Success {
		log.Printf("Response failed: %+v", string(body))
		return fmt.Errorf("failed to create record: %+v", string(body))
	}

	return nil
}

// updateRecord updates DNS A Record with new IP.
func (provider *DNSProvider) updateRecord(record DNSRecord, newIP string) string {
	var r DNSRecordUpdateResponse
	var lastIP string
	record.SetIP(newIP)
	j, _ := json.Marshal(record)
	req, client := provider.newRequest("PUT",
		"/zones/"+record.ZoneID+"/dns_records/"+record.ID,
		bytes.NewBuffer(j),
	)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request error:", err)
		return ""
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalf("Decoder error: %+v", err)
		log.Printf("Response body: %+v", string(body))
		return ""
	}

	if !r.Success {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Response failed: %+v", string(body))
	} else {
		log.Printf("Record updated: %+v - %+v", record.Name, record.IP)
		lastIP = record.IP
	}

	return lastIP
}

// recordTracked checks if record is present in domain conf.
func recordTracked(domain *settings.Domain, record *DNSRecord) bool {
	for _, subDomain := range domain.SubDomains {
		sd := fmt.Sprintf("%s.%s", subDomain, domain.DomainName)
		if record.Name == sd {
			return true
		} else if subDomain == utils.RootDomain && record.Name == domain.DomainName {
			return true
		}
	}

	return false
}
