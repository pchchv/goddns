package linode

import (
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
