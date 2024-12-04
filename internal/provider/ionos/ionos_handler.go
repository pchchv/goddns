package ionos

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const BaseURL = "https://api.hosting.ionos.com/dns/v1/"

type recordResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RootName string `json:"rootName"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Prio     int    `json:"prio"`
	Disabled bool   `json:"disabled"`
}

type zoneResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type recordListResponse struct {
	zoneResponse
	Records []recordResponse `json:"records"`
}

type DNSProvider struct {
	configuration *settings.Settings
	client        *http.Client
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
	provider.client = utils.GetHTTPClient(provider.configuration)
}

func (provider *DNSProvider) getData(endpoint string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-API-Key", provider.configuration.LoginToken)
	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := provider.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get data from " + BaseURL + endpoint + ", status code: " + resp.Status)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (provider *DNSProvider) getZoneID(domainName string) (string, error) {
	body, err := provider.getData("zones", nil)
	if err != nil {
		return "", err
	}

	var zones []zoneResponse
	if err = json.Unmarshal(body, &zones); err != nil {
		return "", err
	}

	for _, zone := range zones {
		if zone.Name == domainName {
			return zone.ID, nil
		}
	}

	return "", errors.New("zone " + domainName + " not found")
}
