package scaleway

// IDFields to filter DNS records for Scaleway API.
type IDFields struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
