package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pchchv/goddns/internal/settings"
)

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

func (c *Controller) GetNetworkSettings(ctx fiber.Ctx) error {
	settings := NetworkSettings{
		IPMode:        c.config.IPType,
		IPUrls:        c.config.IPUrls,
		IPV6Urls:      c.config.IPV6Urls,
		UseProxy:      c.config.UseProxy,
		SkipSSLVerify: c.config.SkipSSLVerify,
		Socks5Proxy:   c.config.Socks5Proxy,
		Webhook:       c.config.Webhook,
		Resolver:      c.config.Resolver,
		IPInterface:   c.config.IPInterface,
	}

	return ctx.JSON(settings)
}
