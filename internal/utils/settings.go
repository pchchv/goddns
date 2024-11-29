package utils

import (
	"errors"

	"github.com/pchchv/goddns/internal/settings"
)

func checkDomains(config *settings.Settings) error {
	for _, d := range config.Domains {
		if d.DomainName == "" {
			return errors.New("domain name should not be empty")
		}

		for _, sd := range d.SubDomains {
			if sd == "" {
				return errors.New("subdomain should not be empty")
			}
		}
	}

	return nil
}
