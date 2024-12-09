package controllers

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/pchchv/goddns/internal/settings"
)

func (c *Controller) GetDomains(ctx fiber.Ctx) error {
	return ctx.JSON(c.config.Domains)
}

func (c *Controller) DeleteDomain(ctx fiber.Ctx) error {
	domainName := ctx.Params("name")
	if domainName == "" {
		return ctx.Status(400).SendString("Domain name is required")
	}

	var domains []settings.Domain
	for _, domain := range c.config.Domains {
		if domain.DomainName != domainName {
			domains = append(domains, domain)
		}
	}

	c.config.Domains = domains
	if err := c.config.SaveSettings(c.configPath); err != nil {
		log.Fatalf("Failed to save settings: %s", err.Error())
		return ctx.Status(500).SendString("Failed to save settings")
	}

	return ctx.JSON(c.config.Domains)
}
