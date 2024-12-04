package he

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const URL = "https://dyn.dns.he.net/nic/update" // API address

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

// updateIP updates subdomain with current IP.
func (provider *DNSProvider) updateIP(domain, subDomain, currentIP string) error {
	values := url.Values{}
	if subDomain != utils.RootDomain {
		values.Add("hostname", fmt.Sprintf("%s.%s", subDomain, domain))
	} else {
		values.Add("hostname", domain)
	}

	values.Add("password", provider.configuration.Password)
	values.Add("myip", currentIP)

	client := utils.GetHTTPClient(provider.configuration)
	req, _ := http.NewRequest("POST", URL, strings.NewReader(values.Encode()))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request error:", err)
		return err
	}

	if body, _ := io.ReadAll(resp.Body); resp.StatusCode == http.StatusOK {
		log.Printf("Update IP success: %s", string(body))
	} else {
		log.Printf("Update IP failed: %s", string(body))
	}

	return nil
}
