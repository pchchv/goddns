package ionos

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func (provider *DNSProvider) getRecord(zoneID, recordName string) (id string, ip string, err error) {
	ipType := utils.IPTypeA
	if provider.configuration.IPType == utils.IPV6 || provider.configuration.IPType == utils.IPTypeAAAA {
		ipType = utils.IPTypeAAAA
	}

	body, err := provider.getData("zones/"+zoneID,
		map[string]string{
			"recordName": recordName,
			"recordType": ipType,
		})
	if err != nil {
		return "", "", err
	}

	var rlp recordListResponse
	if err = json.Unmarshal(body, &rlp); err != nil {
		return "", "", err
	}

	if len(rlp.Records) > 0 {
		return rlp.Records[0].ID, rlp.Records[0].Content, nil
	}

	return "", "", errors.New("record " + recordName + " not found")
}

func (provider *DNSProvider) putData(endpoint string, params map[string]any) (err error) {
	var body []byte
	if params != nil {
		body, err = json.Marshal(params)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(http.MethodPut, BaseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", provider.configuration.LoginToken)
	resp, err := provider.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to PUT " + endpoint + ", status: " + resp.Status)
	}

	defer resp.Body.Close()

	return nil
}

func (provider *DNSProvider) updateRecord(zoneID, recordID, recordName, ip string) (err error) {
	if err = provider.putData(fmt.Sprintf("zones/%s/records/%s", zoneID, recordID), map[string]any{"content": ip}); err != nil {
		return errors.New("failed to update record " + recordName + ": " + err.Error())
	}

	log.Printf("Updated record %s to %s", recordName, ip)

	return nil
}
