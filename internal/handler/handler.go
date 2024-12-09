package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/pchchv/goddns/internal/provider"
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
	"github.com/pchchv/goddns/pkg/ip"
	"github.com/pchchv/goddns/pkg/notification"
	"github.com/pchchv/goddns/pkg/webhook"
)

var (
	errEmptyResult = errors.New("empty result")
	errEmptyDomain = errors.New("NXDOMAIN")
)

type Handler struct {
	ctx                 context.Context
	Configuration       *settings.Settings
	dnsProvider         provider.IDNSProvider
	notificationManager notification.INotificationManager
	ipManager           *ip.IPHelper
	cachedIP            string
}

func (handler *Handler) Init() {
	handler.ipManager.UpdateConfiguration(handler.Configuration)
}

func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
	handler.notificationManager = notification.GetNotificationManager(handler.Configuration)
	handler.ipManager = ip.GetIPHelperInstance(handler.Configuration)
}

func (handler *Handler) SetProvider(provider provider.IDNSProvider) {
	handler.dnsProvider = provider
}

func (handler *Handler) SetContext(ctx context.Context) {
	handler.ctx = ctx
}

func (handler *Handler) UpdateIP(domain *settings.Domain) error {
	ip := handler.ipManager.GetCurrentIP()
	if ip == handler.cachedIP {
		log.Printf("IP (%s) matches cached IP (%s), skipping", ip, handler.cachedIP)
		return nil
	} else if ip == "" {
		if handler.Configuration.RunOnce {
			return errors.New("fail to get current IP")
		}
		return nil
	}

	if err := handler.updateDNS(domain, ip); err != nil {
		if handler.Configuration.RunOnce {
			return errors.New(err.Error() + ": fail to update DNS")
		}
		log.Fatal(err)
		return nil
	}

	handler.cachedIP = ip
	log.Printf("Cached IP address: %s", ip)
	return nil
}

func (handler *Handler) updateDNS(domain *settings.Domain, ip string) error {
	var updatedDomains []string
	for _, subdomainName := range domain.SubDomains {
		var hostname string
		if subdomainName != utils.RootDomain {
			hostname = subdomainName + "." + domain.DomainName
		} else {
			hostname = domain.DomainName
		}

		lastIP, err := utils.ResolveDNS(hostname, handler.Configuration.Resolver, handler.Configuration.IPType)
		if err != nil && (errors.Is(err, errEmptyResult) || errors.Is(err, errEmptyDomain)) {
			log.Fatalf("Failed to resolve DNS for domain: %s, error: %s", hostname, err)
			continue
		}

		// check against the current known IP, if no change, skip update
		if ip == lastIP {
			log.Printf("IP is the same as cached one (%s). Skip update.", ip)
		} else {
			if err := handler.dnsProvider.UpdateIP(domain.DomainName, subdomainName, ip); err != nil {
				return err
			}

			updatedDomains = append(updatedDomains, subdomainName)

			// execute webhook when it is enabled
			if handler.Configuration.Webhook.Enabled {
				if err := webhook.GetWebhook(handler.Configuration).Execute(hostname, ip); err != nil {
					return err
				}
			}
		}
	}

	if len(updatedDomains) > 0 {
		successMessage := fmt.Sprintf("[ %s ] of %s", strings.Join(updatedDomains, ", "), domain.DomainName)
		handler.notificationManager.Send(successMessage, ip)
	}

	return nil
}
