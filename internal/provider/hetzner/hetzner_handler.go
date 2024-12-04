package hetzner

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const BaseURL = "https://dns.hetzner.com/api/v1/" // API address

type Record struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    int64  `json:"ttl"`
	ZoneID string `json:"zone_id"`
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

func (provider *DNSProvider) getData(endpoint string, param string, value string) ([]byte, error) {
	req, _ := http.NewRequest("GET", BaseURL+endpoint, nil)
	q := req.URL.Query()
	q.Add(param, value)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Auth-API-Token", provider.configuration.LoginToken)
	resp, err := provider.client.Do(req)
	if err != nil {
		log.Fatal("Error in fetching: ", err)
		return nil, err
	}

	respBody, _ := io.ReadAll(resp.Body)
	return respBody, nil
}

func (provider *DNSProvider) getZoneID(zoneName string) (string, error) {
	type Zone struct {
		ID string `json:"id"`
	}

	type GetAllZonesResponse struct {
		Zones []Zone `json:"zones"`
	}

	respBody, err := provider.getData("zones", "name", zoneName)
	if err != nil {
		return "", err
	}

	response := GetAllZonesResponse{}
	if err = json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}

	if len(response.Zones) == 0 {
		return "", err
	}

	if len(response.Zones) > 1 {
		return "", err
	}

	return response.Zones[0].ID, nil
}

func (provider *DNSProvider) getRecord(recordName string, zoneID string, Type string) (Record, error) {
	type GetRecordsResult struct {
		Records []Record `json:"records"`
	}

	response := GetRecordsResult{}
	respBody, err := provider.getData("records", "zone_id", zoneID)
	if err != nil {
		return Record{}, err
	}

	if err = json.Unmarshal(respBody, &response); err != nil {
		return Record{}, err
	}

	if len(response.Records) == 0 {
		log.Fatal("Zone doesn't have any records")
		return Record{}, errors.New("zone doesn't have an records")
	}

	outRecord := Record{}
	if Type == "IPv6" {
		Type = utils.IPTypeAAAA
	} else {
		Type = utils.IPTypeA
	}

	found := false
	for _, record := range response.Records {

		if record.Name == recordName && record.Type == Type {
			found = true
			outRecord = record
			break
		}
	}

	if found {
		return outRecord, nil
	}

	return outRecord, errors.New("no record matching value and type found")
}
