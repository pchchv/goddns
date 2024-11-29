package provider

import "github.com/pchchv/goddns/internal/settings"

type IDNSProvider interface {
	Init(conf *settings.Settings)
	UpdateIP(domainName, subdomainName, ip string) error
}
