package infomaniak

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const URL = "https://%s:%s@infomaniak.com/nic/update?hostname=%s.%s&myip=%s" // API address

type DNSProvider struct {
	configuration *settings.Settings
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	return provider.updateIP(domainName, subdomainName, ip)
}

// updateIP update subdomain with current IP.
func (provider *DNSProvider) updateIP(domain, subDomain, currentIP string) error {
	client := utils.GetHTTPClient(provider.configuration)
	resp, err := client.Get(fmt.Sprintf(URL,
		provider.configuration.Email,
		provider.configuration.Password,
		subDomain,
		domain,
		currentIP))
	if err != nil {
		log.Fatal("Failed to update sub domain:", subDomain)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Update IP failed: %s", string(body))
		return fmt.Errorf("update IP failed: %s", string(body))
	}

	if strings.Contains(string(body), "good") {
		log.Printf("Update IP success: %s", string(body))
	} else if strings.Contains(string(body), "nochg") {
		log.Printf("IP not changed: %s", string(body))
	} else {
		return fmt.Errorf("update IP failed: %s", string(body))
	}

	return nil
}
