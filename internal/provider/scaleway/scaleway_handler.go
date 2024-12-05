package scaleway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const URL = "https://api.scaleway.com/domain/v2beta1/dns-zones/%s/records"

// IDFields to filter DNS records for Scaleway API.
type IDFields struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Record struct {
	Name    string `json:"name"`
	Data    string `json:"data"`
	TTL     int    `json:"ttl"`
	Comment string `json:"comment"`
}

type SetRecord struct {
	IDFields IDFields `json:"id_fields"`
	Records  []Record `json:"records"`
}

type DNSChange struct {
	Set SetRecord `json:"set"`
}

type DNSUpdateRequest struct {
	Changes []DNSChange `json:"changes"`
}

type DNSProvider struct {
	configuration *settings.Settings
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

func (provider *DNSProvider) UpdateIP(domainName string, subdomainName string, ip string) error {
	log.Printf("%s.%s - Start to update record IP...", subdomainName, domainName)
	if err := provider.updateIP(domainName, subdomainName, ip); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (provider *DNSProvider) getRecordType() (string, error) {
	if strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		return utils.IPTypeA, nil
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		return utils.IPTypeAAAA, nil
	}

	return "", errors.New("must specify \"ip_type\" in config for Scaleway")
}

// updateIP update subdomain with current IP.
func (provider *DNSProvider) updateIP(domain, subDomain, currentIP string) error {
	recordType, err := provider.getRecordType()
	if err != nil {
		return err
	}

	reqBody := DNSUpdateRequest{Changes: []DNSChange{{SetRecord{
		IDFields: IDFields{
			Name: subDomain,
			Type: recordType,
		},
		Records: []Record{
			{
				Name:    subDomain,
				Data:    currentIP,
				TTL:     provider.configuration.Interval,
				Comment: "Set by GoDNS",
			},
		},
	}}}}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return errors.New("failed to encode request body as json")
	}

	req, _ := http.NewRequest("PATCH", fmt.Sprintf(URL, domain), bytes.NewReader(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", provider.configuration.LoginToken)
	if provider.configuration.UserAgent != "" {
		req.Header.Add("User-Agent", provider.configuration.UserAgent)
	}

	client := utils.GetHTTPClient(provider.configuration)
	log.Printf("Requesting update for '%s.%s': '%v'", subDomain, domain, reqBody)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request error:", err)
		return errors.New("failed to complete update request")
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Update failed for '%s.%s': %s", subDomain, domain, string(body))
		return errors.New("update IP failed with status " + string(resp.StatusCode))
	}

	log.Printf("Update IP success for '%s.%s': '%s'", subDomain, domain, string(body))

	return nil
}
