package ovh

import (
	"fmt"
	"log"
	"strings"

	"github.com/ovh/go-ovh/ovh"
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

type Record struct {
	Zone      string `json:"zone"`
	TTL       int    `json:"ttl"`
	Value     string `json:"target"`
	SubDomain string `json:"subDomain"`
	Type      string `json:"fieldType"`
	ID        int    `json:"id"`
}

type DNSProvider struct {
	configuration *settings.Settings
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

func (provider *DNSProvider) UpdateIP(domainName string, subdomainName string, ip string) error {
	client, err := ovh.NewClient(
		"ovh-eu",
		provider.configuration.AppKey,
		provider.configuration.AppSecret,
		provider.configuration.ConsumerKey,
	)
	if err != nil {
		log.Fatal("OVH Client error: ", err)
		return err
	}

	var IDs []int
	query := fmt.Sprintf("/domain/zone/%s/record?subDomain=%s", domainName, subdomainName)
	if err = client.Get(query, &IDs); err != nil {
		log.Fatal("Fetch error")
		return err
	}

	if len(IDs) < 1 {
		log.Fatal("No matching records")
		return fmt.Errorf("no matching records")
	}

	outrec := Record{}
	for _, id := range IDs {
		record := Record{}
		if err = client.Get(fmt.Sprintf("/domain/zone/%s/record/%s", domainName, fmt.Sprint(id)), &record); err != nil {
			log.Fatal("Fetch error on get record: ", id)
			return err
		}

		if strings.ToUpper(provider.configuration.IPType) == provider.recordTypeToIPType(record.Type) {
			outrec = record
			break
		}
	}

	if outrec.ID == 0 {
		log.Fatal("No fitting record type found")
		return fmt.Errorf("no fitting record type found")
	}

	// update IP
	outrec.Value = ip
	if err = client.Put(fmt.Sprintf("/domain/zone/%s/record/%s", domainName, fmt.Sprint(outrec.ID)), outrec, nil); err != nil {
		log.Fatal("Error while Updating record: ", outrec.SubDomain, outrec.Zone)
		return err
	}

	// refresh zone
	if err = client.Post(fmt.Sprintf("/domain/zone/%s/refresh", domainName), nil, nil); err != nil {
		log.Fatal("Applying new records failed")
		return err
	}

	return nil
}

func (provider *DNSProvider) recordTypeToIPType(Type string) string {
	if Type == utils.IPTypeAAAA {
		return utils.IPV6
	}

	return utils.IPV4
}
