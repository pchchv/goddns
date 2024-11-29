package alidns

import "sync"

var (
	instance *AliDNS
	once     sync.Once
)

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

type domainRecordsResp struct {
	RequestID     string `json:"RequestId"`
	TotalCount    int
	PageNumber    int
	PageSize      int
	DomainRecords domainRecords
}

// AliDNS token.
type AliDNS struct {
	AccessKeyID     string
	AccessKeySecret string
	IPType          string
}

// NewAliDNS function creates instance of AliDNS and return.
func NewAliDNS(key, secret, ipType string) *AliDNS {
	once.Do(func() {
		instance = &AliDNS{
			AccessKeyID:     key,
			AccessKeySecret: secret,
			IPType:          ipType,
		}
	})
	return instance
}
