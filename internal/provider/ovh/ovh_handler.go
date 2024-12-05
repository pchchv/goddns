package ovh

type Record struct {
	Zone      string `json:"zone"`
	TTL       int    `json:"ttl"`
	Value     string `json:"target"`
	SubDomain string `json:"subDomain"`
	Type      string `json:"fieldType"`
	ID        int    `json:"id"`
}
