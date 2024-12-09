package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
	"github.com/pchchv/goddns/pkg/ip"
)

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

func (c *Controller) GetSubDomains() (count int) {
	// get the total number of all the sub domains
	for _, domain := range c.config.Domains {
		count += len(domain.SubDomains)
	}
	return
}

func (c *Controller) GetBasicInfo(ctx fiber.Ctx) error {
	return ctx.JSON(BasicInfo{
		Version:      utils.Version,
		StartTime:    utils.StartTime,
		DomainNum:    c.getDomains(),
		SubDomainNum: c.GetSubDomains(),
		Domains:      c.config.Domains,
		PublicIP:     ip.GetIPHelperInstance(c.config).GetCurrentIP(),
		IPMode:       strings.ToUpper(c.config.IPType),
		Provider:     c.config.Provider,
	})
}

func (c *Controller) getDomains() int {
	// count the total number of domains
	return len(c.config.Domains)
}
