// Package resolver is a simple dns resolver
// based on miekg/dns
package resolver

import (
	"errors"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/miekg/dns"
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

// NewFromResolvConf initializes DnsResolver from resolv.conf like file.
func NewFromResolvConf(path string) (*DNSResolver, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &DNSResolver{}, errors.New("no such file or directory: " + path)
	}

	config, err := dns.ClientConfigFromFile(path)
	if err != nil {
		return &DNSResolver{}, err
	}

	var servers []string
	for _, ipAddress := range config.Servers {
		servers = append(servers, net.JoinHostPort(ipAddress, "53"))
	}

	return &DNSResolver{servers, len(servers) * 2, rand.New(rand.NewSource(time.Now().UnixNano()))}, err
}
