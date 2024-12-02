package digitalocean

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDNSResponseToJSON(t *testing.T) {
	s := strings.NewReader(`
    {
        "domain_records": [
            {
                "id": 12345678,
                "type": "A",
                "name": "potato",
                "data": "127.0.0.1",
                "priority": null,
                "port": null,
                "ttl": 3600,
                "weight": null,
                "flags": null,
                "tag": null
            }
        ],
        "links": {},
        "meta": {
            "total": 1
        }
    }`)

	var resp DomainRecordsResponse
	if err := json.NewDecoder(s).Decode(&resp); err != nil {
		t.Error(err.Error())
	}

	if resp.Records[0].ID != 12345678 {
		t.Errorf("ID Error: %#v != 12345678 ", resp.Records[0].ID)
	}

	if resp.Records[0].Name != "potato" {
		t.Errorf("Name Error: %#v != potato", resp.Records[0].Name)
	}
}
