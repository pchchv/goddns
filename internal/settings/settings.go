package settings

type Domain struct {
	DomainName string   `json:"domain_name" yaml:"domain_name"`
	SubDomains []string `json:"sub_domains" yaml:"sub_domains"`
}
