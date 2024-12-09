package controllers

import "github.com/pchchv/goddns/internal/settings"

type NetworkSettings struct {
	IPMode        string           `json:"ip_mode"`
	IPUrls        []string         `json:"ip_urls"`
	IPV6Urls      []string         `json:"ipv6_urls"`
	UseProxy      bool             `json:"use_proxy"`
	SkipSSLVerify bool             `json:"skip_ssl_verify"`
	Socks5Proxy   string           `json:"socks5_proxy"`
	Webhook       settings.Webhook `json:"webhook,omitempty"`
	Resolver      string           `json:"resolver"`
	IPInterface   string           `json:"ip_interface"`
}
