package provider

import (
	"errors"

	"github.com/pchchv/goddns/internal/provider/alidns"
	"github.com/pchchv/goddns/internal/provider/cloudflare"
	"github.com/pchchv/goddns/internal/provider/digitalocean"
	"github.com/pchchv/goddns/internal/provider/dnspod"
	"github.com/pchchv/goddns/internal/provider/dreamhost"
	"github.com/pchchv/goddns/internal/provider/duck"
	"github.com/pchchv/goddns/internal/provider/dynu"
	"github.com/pchchv/goddns/internal/provider/dynv6"
	"github.com/pchchv/goddns/internal/provider/google"
	"github.com/pchchv/goddns/internal/provider/he"
	"github.com/pchchv/goddns/internal/provider/hetzner"
	"github.com/pchchv/goddns/internal/provider/infomaniak"
	"github.com/pchchv/goddns/internal/provider/ionos"
	"github.com/pchchv/goddns/internal/provider/linode"
	"github.com/pchchv/goddns/internal/provider/loopiase"
	"github.com/pchchv/goddns/internal/provider/noip"
	"github.com/pchchv/goddns/internal/provider/ovh"
	"github.com/pchchv/goddns/internal/provider/scaleway"
	"github.com/pchchv/goddns/internal/provider/strato"
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

func GetProvider(conf *settings.Settings) (provider IDNSProvider, err error) {
	switch conf.Provider {
	case utils.CLOUDFLARE:
		provider = &cloudflare.DNSProvider{}
	case utils.DIGITALOCEAN:
		provider = &digitalocean.DNSProvider{}
	case utils.DNSPOD:
		provider = &dnspod.DNSProvider{}
	case utils.DREAMHOST:
		provider = &dreamhost.DNSProvider{}
	case utils.HE:
		provider = &he.DNSProvider{}
	case utils.ALIDNS:
		provider = &alidns.DNSProvider{}
	case utils.GOOGLE:
		provider = &google.DNSProvider{}
	case utils.DUCK:
		provider = &duck.DNSProvider{}
	case utils.NOIP:
		provider = &noip.DNSProvider{}
	case utils.SCALEWAY:
		provider = &scaleway.DNSProvider{}
	case utils.DYNV6:
		provider = &dynv6.DNSProvider{}
	case utils.LINODE:
		provider = &linode.DNSProvider{}
	case utils.STRATO:
		provider = &strato.DNSProvider{}
	case utils.LOOPIASE:
		provider = &loopiase.DNSProvider{}
	case utils.INFOMANIAK:
		provider = &infomaniak.DNSProvider{}
	case utils.HETZNER:
		provider = &hetzner.DNSProvider{}
	case utils.OVH:
		provider = &ovh.DNSProvider{}
	case utils.DYNU:
		provider = &dynu.DNSProvider{}
	case utils.IONOS:
		provider = &ionos.DNSProvider{}
	}

	if provider != nil {
		provider.Init(conf)
	} else {
		err = errors.New("Unknown provider " + conf.Provider)
	}

	return
}
