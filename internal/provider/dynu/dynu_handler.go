package dynu

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const URL = "https://api.dynu.com/nic/update?hostname=%s&password=%s&%s" // API address for dynu

type DNSProvider struct {
	configuration *settings.Settings
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	hostname := subdomainName + "." + domainName
	client := utils.GetHTTPClient(provider.configuration)
	return provider.update(client, hostname, subdomainName, ip)
}

func (provider *DNSProvider) update(client *http.Client, hostname, subdomain string, currentIP string) error {
	var ip string
	if strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		ip = fmt.Sprintf("myip=%s", currentIP)
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		ip = fmt.Sprintf("myipv6=%s", currentIP)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf(
		URL,
		hostname,
		utils.GetMD5Hash(provider.configuration.Password),
		ip), nil)
	if provider.configuration.UserAgent != "" {
		req.Header.Add("User-Agent", provider.configuration.UserAgent)
	}

	// update IP with HTTP GET request
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		log.Fatal("Failed to update sub domain:", subdomain)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if body, err := io.ReadAll(resp.Body); err != nil || !strings.Contains(string(body), "good") {
		log.Fatal("Failed to update the IP", err)
		return err
	}

	log.Printf("IP updated to: %s", currentIP)

	return nil
}
