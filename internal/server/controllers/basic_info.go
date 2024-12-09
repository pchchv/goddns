package controllers

import "github.com/pchchv/goddns/internal/settings"

type BasicInfo struct {
	Version      string            `json:"version"`
	StartTime    int64             `json:"start_time"`
	DomainNum    int               `json:"domain_num"`
	SubDomainNum int               `json:"sub_domain_num"`
	Domains      []settings.Domain `json:"domains"`
	PublicIP     string            `json:"public_ip"`
	IPMode       string            `json:"ip_mode"`
	Provider     string            `json:"provider"`
}

func (c *Controller) getDomains() int {
	// count the total number of domains
	return len(c.config.Domains)
}
