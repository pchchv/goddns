package utils

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
	IPTypeA        = "A"
	IPTypeAAAA     = "AAAA"
	IPV4           = "IPV4"
	IPV6           = "IPV6"
	LINODE         = "Linode"
	LOOPIASE       = "LoopiaSE"
	NOIP           = "NoIP"
	OVH            = "OVH"
	SCALEWAY       = "Scaleway"
	STRATO         = "Strato"
	RootDomain     = "@"
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

// Version is current version of GoDNS.
var Version = "v0.1"
