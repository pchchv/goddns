package utils

import (
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/pchchv/goddns/pkg/resolver"
)

// ResolveDNS will query DNS for a given hostname.
func ResolveDNS(hostname, r, ipType string) (string, error) {
	var dnsType uint16
	if ipType == "" || strings.ToUpper(ipType) == IPV4 {
		dnsType = dns.TypeA
	} else {
		dnsType = dns.TypeAAAA
	}

	// if no DNS server is set in config file,
	// falls back to default resolver
	if r == "" {
		dnsAddress, err := net.LookupHost(hostname)
		if err != nil {
			return "<nil>", err
		}

		return dnsAddress[0], nil
	}

	res := resolver.New([]string{r})
	// in case of i/o timeout
	res.RetryTimes = 5
	ip, err := res.LookupHost(hostname, dnsType)
	if err != nil {
		return "<nil>", err
	}

	return ip[0].String(), nil
}
