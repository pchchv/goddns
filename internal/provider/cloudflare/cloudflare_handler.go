package cloudflare

import "github.com/pchchv/goddns/internal/settings"

// Zone object with id and name.
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ZoneResponse is a wrapper for Zones.
type ZoneResponse struct {
	Zones   []Zone `json:"result"`
	Success bool   `json:"success"`
}

// DNSRecord for Cloudflare API.
type DNSRecord struct {
	ID      string `json:"id"`
	IP      string `json:"content"`
	Name    string `json:"name"`
	Proxied bool   `json:"proxied"`
	Type    string `json:"type"`
	ZoneID  string `json:"zone_id"`
	TTL     int32  `json:"ttl"`
}

// SetIP updates DNSRecord.IP.
func (r *DNSRecord) SetIP(ip string) {
	r.IP = ip
}

type DNSRecordUpdateResponse struct {
	Record  DNSRecord `json:"result"`
	Success bool      `json:"success"`
}

type DNSRecordResponse struct {
	Records []DNSRecord `json:"result"`
	Success bool        `json:"success"`
}

type DNSProvider struct {
	configuration *settings.Settings
	API           string
}
