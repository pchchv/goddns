package digitalocean

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

func (provider *DNSProvider) getCurrentDomain(domainName string) *settings.Domain {
	for _, domain := range provider.configuration.Domains {
		domain := domain
		if domain.DomainName == domainName {
			return &domain
		}
	}

	return nil
}

// updateRecord updates DNS Record with new IP.
func (provider *DNSProvider) updateRecord(domainName string, record DNSRecord, newIP string) string {
	record.SetIP(newIP)
	j, _ := json.Marshal(record)
	req, client := provider.newRequest("PUT",
		fmt.Sprintf("/domains/%s/records/%d", domainName, record.ID),
		bytes.NewBuffer(j),
	)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request error:", err)
		return ""
	}
	defer resp.Body.Close()

	var r DNSRecord
	body, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &r); err != nil {
		log.Fatalf("Decoder error: %+v", err)
		log.Printf("Response body: %+v", string(body))
		return ""
	}

	log.Printf("Record updated: %+v - %+v", record.Name, record.IP)
	return record.IP
}

func (provider *DNSProvider) createRecord(domain, subDomain, ip string) error {
	recordType := provider.getRecordType()
	newRecord := DNSRecord{
		Type: recordType,
		IP:   ip,
		TTL:  int32(provider.configuration.Interval),
	}

	if subDomain == utils.RootDomain {
		newRecord.Name = utils.RootDomain
	} else {
		newRecord.Name = subDomain
	}

	content, err := json.Marshal(newRecord)
	if err != nil {
		log.Fatalf("Encoder error: %+v", err)
		return err
	}

	req, client := provider.newRequest("POST", fmt.Sprintf("/domains/%s/records", domain), bytes.NewBuffer(content))
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

	var r DNSRecord
	if err = json.Unmarshal(body, &r); err != nil {
		log.Fatalf("Response decoder error: %+v", err)
		log.Printf("Response body: %+v", string(body))
		return err
	}

	return nil
}

// recordTracked checks if record is present in domain conf.
func recordTracked(domain *settings.Domain, record *DNSRecord) bool {
	for _, subDomain := range domain.SubDomains {
		if record.Name == subDomain {
			return true
		}
	}

	return false
}
