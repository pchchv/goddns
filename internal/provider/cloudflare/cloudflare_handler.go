package cloudflare

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
