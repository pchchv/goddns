// Package resolver is a simple dns resolver
// based on miekg/dns
package resolver

import "math/rand"

type DNSResolver struct {
	Servers    []string
	RetryTimes int
	r          *rand.Rand
}
