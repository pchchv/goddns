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
