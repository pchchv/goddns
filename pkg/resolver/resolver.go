// Package resolver is a simple dns resolver
// based on miekg/dns
package resolver

import (
	"math/rand"
	"net"
	"time"
)

type DNSResolver struct {
	Servers    []string
	RetryTimes int
	r          *rand.Rand
}

// New initializes DnsResolver.
func New(servers []string) *DNSResolver {
	for i := range servers {
		servers[i] = net.JoinHostPort(servers[i], "53")
	}

	return &DNSResolver{servers, len(servers) * 2, rand.New(rand.NewSource(time.Now().UnixNano()))}
}
