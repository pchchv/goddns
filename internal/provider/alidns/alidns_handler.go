package alidns

import (
	"fmt"
	"log"
)

type DNSProvider struct {
	aliDNS *AliDNS
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	log.Printf("%s.%s - Start to update record IP...", subdomainName, domainName)
	records := provider.aliDNS.GetDomainRecords(domainName, subdomainName)
	if len(records) == 0 {
		log.Fatalf("Cannot get subdomain [%s] from AliDNS.", subdomainName)
		return fmt.Errorf("cannot get subdomain [%s] from AliDNS", subdomainName)
	}

	if records[0].Value != ip {
		records[0].Value = ip
		if err := provider.aliDNS.UpdateDomainRecord(records[0]); err != nil {
			return fmt.Errorf("failed to update IP for subdomain: %s", subdomainName)
		}
		log.Printf("IP updated for subdomain: %s", subdomainName)
	} else {
		log.Printf("IP not changed for subdomain: %s", subdomainName)
	}

	return nil
}
