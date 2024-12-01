package cloudflare

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResponseToJSON(t *testing.T) {
	s := strings.NewReader(`
    {
        "errors": [],
        "messages": [],
        "result": [
            {
                "id": "mk2b6fa491c12445a4376666a32429e1",
                "name": "example.com",
                "status": "active"
            }
        ],
        "result_info": {
            "count": 1,
            "page": 1,
            "per_page": 20,
            "total_count": 1,
            "total_pages": 1
        },
        "success": true
    }`)

	var resp ZoneResponse
	if err := json.NewDecoder(s).Decode(&resp); err != nil {
		t.Error(err.Error())
	}

	if resp.Success != true {
		t.Errorf("Success Error: %#v != true ", resp.Success)
	}

	if resp.Zones[0].ID != "mk2b6fa491c12445a4376666a32429e1" {
		t.Errorf("ID Error: %#v != mk2b6fa491c12445a4376666a32429e1 ", resp.Zones[0].ID)
	}

	if resp.Zones[0].Name != "example.com" {
		t.Errorf("Name Error: %#v != example.com", resp.Zones[0].Name)
	}
}

func TestDNSResponseToJSON(t *testing.T) {
	s := strings.NewReader(`
    {
        "errors": [],
        "messages": [],
        "result": [
            {
                "content": "127.0.0.1",
                "id": "F11cc63e02a42d38174b8e7c548a7b6f",
                "name": "example.com",
                "type": "A",
                "zone_id": "mk2b6fa491c12445a4376666a32429e1",
                "zone_name": "example.com"
            }
        ],
        "success": true
    }`)

	var resp DNSRecordResponse
	if err := json.NewDecoder(s).Decode(&resp); err != nil {
		t.Error(err.Error())
	}

	if resp.Success != true {
		t.Errorf("Success Error: %#v != true ", resp.Success)
	}

	if resp.Records[0].ID != "F11cc63e02a42d38174b8e7c548a7b6f" {
		t.Errorf("ID Error: %#v != F11cc63e02a42d38174b8e7c548a7b6f ", resp.Records[0].ID)
	}

	if resp.Records[0].Name != "example.com" {
		t.Errorf("Name Error: %#v != example.com", resp.Records[0].Name)
	}
}
