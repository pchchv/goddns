package linode

import (
	"context"
	"errors"
	"log"

	"github.com/linode/linodego"
	"github.com/pchchv/goddns/internal/settings"
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
