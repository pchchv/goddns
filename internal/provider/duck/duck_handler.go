package duck

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const URL = "https://www.duckdns.org/update?domains=%s&token=%s&%s" // API address for Duck DNS

type DNSProvider struct {
	configuration *settings.Settings
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

func (provider *DNSProvider) updateIP(domainName, subdomainName, currentIP string) error {
	var ip string
	if strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		ip = fmt.Sprintf("ip=%s", currentIP)
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		ip = fmt.Sprintf("ipv6=%s", currentIP)
	}

	client := utils.GetHTTPClient(provider.configuration)
	// update IP with HTTP GET request
	resp, err := client.Get(fmt.Sprintf(URL, subdomainName, provider.configuration.LoginToken, ip))
	if err != nil {
		// handle error
		log.Fatalf("Failed to update sub domain: %s.%s, error: %s", domainName, subdomainName, err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if body, err := io.ReadAll(resp.Body); err != nil || string(body) != "OK" {
		log.Fatalf("Failed to update the IP, error: %s, body: %s", err, string(body))
		return err
	}

	log.Printf("IP updated to: %s", ip)

	return nil
}
