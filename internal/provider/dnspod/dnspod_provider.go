package dnspod

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const providerURL = "https://dnsapi.cn"

type DNSProvider struct {
	configuration *settings.Settings
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

// generateHeader generates the request header for DNSPod API.
func (provider *DNSProvider) generateHeader(content url.Values) url.Values {
	header := url.Values{}
	if provider.configuration.LoginToken != "" {
		header.Add("login_token", provider.configuration.LoginToken)
	}

	header.Add("format", "json")
	header.Add("lang", "en")
	header.Add("error_on_empty", "no")

	for k := range content {
		header.Add(k, content.Get(k))
	}

	return header
}

// postData post data and invoke DNSPod API.
func (provider *DNSProvider) postData(url string, content url.Values) (string, error) {
	client := utils.GetHTTPClient(provider.configuration)

	if client == nil {
		return "", errors.New("failed to create HTTP client")
	}

	values := provider.generateHeader(content)
	req, _ := http.NewRequest("POST", providerURL+url, strings.NewReader(values.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("GoDNS/0.1 (%s)", ""))

	response, err := client.Do(req)
	if err != nil {
		log.Fatal("Post failed:", err)
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Failed to close body:", err)
		}
	}(response.Body)

	resp, _ := io.ReadAll(response.Body)

	return string(resp), nil
}
