package scaleway

// IDFields to filter DNS records for Scaleway API.
type IDFields struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Record struct {
	Name    string `json:"name"`
	Data    string `json:"data"`
	TTL     int    `json:"ttl"`
	Comment string `json:"comment"`
}

type SetRecord struct {
	IDFields IDFields `json:"id_fields"`
	Records  []Record `json:"records"`
}

type DNSChange struct {
	Set SetRecord `json:"set"`
}

type DNSUpdateRequest struct {
	Changes []DNSChange `json:"changes"`
}
