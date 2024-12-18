package linode

import (
	"context"
	"errors"
	"log"

	"github.com/linode/linodego"
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

type DNSProvider struct {
	linodeClient *linodego.Client
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	httpClient, err := CreateHTTPClient(conf)
	if err != nil {
		log.Panic(err)
	}

	linodeAPIClient := linodego.NewClient(httpClient)
	linodeAPIClient.SetDebug(conf.DebugInfo)
	provider.linodeClient = &linodeAPIClient
}

func (provider *DNSProvider) UpdateIP(domain, subdomain, ip string) error {
	if subdomain == utils.RootDomain {
		subdomain = ""
	}

	domainID, err := provider.getDomainID(domain)
	if err != nil {
		return err
	}

	recordExists, recordID, err := provider.getDomainRecordID(domainID, subdomain)
	if err != nil {
		return err
	} else if !recordExists {
		recordID, _ = provider.createDomainRecord(domainID, subdomain)
	}

	return provider.updateDomainRecord(domainID, recordID, ip)
}

func (provider *DNSProvider) getDomainID(name string) (int, error) {
	f := linodego.Filter{}
	f.AddField(linodego.Eq, "domain", name)
	fStr, err := f.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	opts := linodego.NewListOptions(0, string(fStr))
	res, err := provider.linodeClient.ListDomains(context.Background(), opts)
	if err != nil {
		return 0, err
	}

	if len(res) == 0 {
		return 0, errors.New("No domains found for name " + name)
	}

	return res[0].ID, nil
}

func (provider *DNSProvider) getDomainRecordID(domainID int, name string) (bool, int, error) {
	res, err := provider.linodeClient.ListDomainRecords(context.Background(), domainID, nil)
	if err != nil {
		return false, 0, err
	} else if len(res) == 0 {
		return false, 0, nil
	}

	for _, record := range res {
		if record.Name == name {
			return true, record.ID, nil
		}
	}

	return false, 0, nil
}

func (provider *DNSProvider) createDomainRecord(domainID int, name string) (int, error) {
	opts := &linodego.DomainRecordCreateOptions{
		Type:   "A",
		Name:   name,
		Target: "127.0.0.1",
		TTLSec: 30,
	}
	record, err := provider.linodeClient.CreateDomainRecord(context.Background(), domainID, *opts)
	if err != nil {
		return 0, err
	}

	return record.ID, nil
}

func (provider *DNSProvider) updateDomainRecord(domainID int, id int, ip string) error {
	opts := &linodego.DomainRecordUpdateOptions{Target: ip}
	_, err := provider.linodeClient.UpdateDomainRecord(context.Background(), domainID, id, *opts)
	return err
}
