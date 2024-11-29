package alidns

type DomainRecord struct {
	DomainName string
	RecordID   string `json:"RecordId"`
	RR         string
	Type       string
	Value      string
	Line       string
	Priority   int
	TTL        int
	Status     string
	Locked     bool
}

type domainRecords struct {
	Record []DomainRecord
}

// AliDNS token.
type AliDNS struct {
	AccessKeyID     string
	AccessKeySecret string
	IPType          string
}
