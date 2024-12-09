package utils

import "time"

const (
	ALIDNS         = "AliDNS"
	CLOUDFLARE     = "Cloudflare"
	DefaultTimeout = 10 // in seconds
	DIGITALOCEAN   = "DigitalOcean"
	DNSPOD         = "DNSPod" // dnspod.cn
	DREAMHOST      = "Dreamhost"
	DUCK           = "DuckDNS"
	DYNU           = "Dynu"
	DYNV6          = "Dynv6"
	GOOGLE         = "Google"
	HE             = "HE" // he.net
	HETZNER        = "Hetzner"
	INFOMANIAK     = "Infomaniak"
	IONOS          = "IONOS"
	IPPattern      = "(" + IPv4Pattern + ")|(" + IPv6Pattern + ")" // regex pattern to match IPV4 and IPV6 address.
	IPTypeA        = "A"
	IPTypeAAAA     = "AAAA"
	IPV4           = "IPV4"
	IPv4Pattern    = `((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)` // regex pattern to match IPV4 address
	IPV6           = "IPV6"
	IPv6Pattern    = `(([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|` + // regex pattern to match IPV6 address
		`(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|` +
		`(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|` +
		`(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))`
	LINODE     = "Linode"
	LOOPIASE   = "LoopiaSE"
	NOIP       = "NoIP"
	OVH        = "OVH"
	SCALEWAY   = "Scaleway"
	STRATO     = "Strato"
	RootDomain = "@"
)

var (
	StartTime = time.Now().Unix()
	Version   = "v0.1" // current version of GoDDNS
)

type ProviderSetting struct {
	Name        string `json:"name" yaml:"name"`
	Username    bool   `json:"username" yaml:"username"`
	Email       bool   `json:"email" yaml:"email"`
	Password    bool   `json:"password" yaml:"password"`
	LoginToken  bool   `json:"login_token" yaml:"login_token"`
	AppKey      bool   `json:"app_key" yaml:"app_key"`
	AppSecret   bool   `json:"app_secret" yaml:"app_secret"`
	ConsumerKey bool   `json:"consumer_key" yaml:"consumer_key"`
}
