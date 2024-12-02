package dreamhost

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const URL = "https://api.dreamhost.com"

type DNSProvider struct {
	configuration *settings.Settings
}

// Init pass DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

// updateDNS can add or remove DNS records.
func (provider *DNSProvider) updateDNS(dns, ip, hostname, action string) error {
	var ipType string
	if provider.configuration.IPType == "" || strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		ipType = utils.IPTypeA
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		ipType = utils.IPTypeAAAA
	}

	// Generates UUID
	uid, _ := uuid.NewRandom()
	values := url.Values{}
	values.Add("record", hostname)
	values.Add("key", provider.configuration.LoginToken)
	values.Add("type", ipType)
	values.Add("unique_id", uid.String())
	switch action {
	case "remove":
		// Build URL query (remove)
		values.Add("cmd", "dns-remove_record")
		values.Add("value", dns)
	case "add":
		// Build URL query (add)
		values.Add("cmd", "dns-add_record")
		values.Add("value", ip)
	default:
		log.Fatalf("Unknown action: %s", action)
		return fmt.Errorf("unknown action: %s", action)
	}

	client := utils.GetHTTPClient(provider.configuration)
	req, _ := http.NewRequest("POST", URL, strings.NewReader(values.Encode()))
	req.SetBasicAuth(provider.configuration.Email, provider.configuration.Password)
	if provider.configuration.UserAgent != "" {
		req.Header.Add("User-Agent", provider.configuration.UserAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request err:", err.Error())
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Failed to close the request body:", err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Update IP failed: %s", string(body))
		return fmt.Errorf("update IP failed: %s", string(body))
	}

	log.Printf("Update IP success: %s", string(body))
	return nil
}
