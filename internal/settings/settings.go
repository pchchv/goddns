package settings

type Domain struct {
	DomainName string   `json:"domain_name" yaml:"domain_name"`
	SubDomains []string `json:"sub_domains" yaml:"sub_domains"`
}

type Webhook struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	URL         string `json:"url" yaml:"url"`
	RequestBody string `json:"request_body" yaml:"request_body"`
}

type WebPanel struct {
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	Addr     string `json:"addr" yaml:"addr"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Mikrotik struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Addr      string `json:"addr" yaml:"addr"`
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	Interface string `json:"interface" yaml:"interface"`
}

type TelegramNotify struct {
	Enabled       bool   `json:"enabled" yaml:"enabled"`
	BotAPIKey     string `json:"bot_api_key" yaml:"bot_api_key"`
	BotAPIKeyFile string `json:"bot_api_key_file" yaml:"bot_api_key_file"`
	ChatID        string `json:"chat_id" yaml:"chat_id"`
	MsgTemplate   string `json:"message_template" yaml:"message_template"`
}

type DiscordNotify struct {
	Enabled         bool   `json:"enabled" yaml:"enabled"`
	BotAPIToken     string `json:"bot_api_token" yaml:"bot_api_token"`
	BotAPITokenFile string `json:"bot_api_token_file" yaml:"bot_api_token_file"`
	Channel         string `json:"channel" yaml:"channel"`
	MsgTemplate     string `json:"message_template" yaml:"message_template"`
}

type SlackNotify struct {
	Enabled         bool   `json:"enabled" yaml:"enabled"`
	BotAPIToken     string `json:"bot_api_token" yaml:"bot_api_token"`
	BotAPITokenFile string `json:"bot_api_token_file" yaml:"bot_api_token_file"`
	Channel         string `json:"channel" yaml:"channel"`
	MsgTemplate     string `json:"message_template" yaml:"message_template"`
}

type PushoverNotify struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Token       string `json:"token" yaml:"token"`
	TokenFile   string `json:"token_file" yaml:"token:_file"`
	User        string `json:"user" yaml:"user"`
	MsgTemplate string `json:"message_template" yaml:"message_template"`
	Device      string `json:"device" yaml:"device"`
	Title       string `json:"title" yaml:"title"`
	Priority    int    `json:"priority" yaml:"priority"`
	HTML        int    `json:"html" yaml:"html"`
}

type MailNotify struct {
	Enabled          bool   `json:"enabled" yaml:"enabled"`
	SMTPServer       string `json:"smtp_server" yaml:"smtp_server"`
	SMTPUsername     string `json:"smtp_username" yaml:"smtp_username"`
	SMTPPassword     string `json:"smtp_password" yaml:"smtp_password"`
	SMTPPasswordFile string `json:"smtp_password_file" yaml:"smtp_password_file"`
	SMTPPort         int    `json:"smtp_port" yaml:"smtp_port"`
	SendFrom         string `json:"send_from" yaml:"send_from"`
	SendTo           string `json:"send_to" yaml:"send_to"`
}

type Notify struct {
	Telegram TelegramNotify `json:"telegram" yaml:"telegram"`
	Mail     MailNotify     `json:"mail" yaml:"mail"`
	Slack    SlackNotify    `json:"slack" yaml:"slack"`
	Discord  DiscordNotify  `json:"discord" yaml:"discord"`
	Pushover PushoverNotify `json:"pushover" yaml:"pushover"`
}

type Settings struct {
	Provider       string   `json:"provider" yaml:"provider"`
	Email          string   `json:"email" yaml:"email"`
	Password       string   `json:"password" yaml:"password"`
	PasswordFile   string   `json:"password_file" yaml:"password_file"`
	LoginToken     string   `json:"login_token" yaml:"login_token"`
	LoginTokenFile string   `json:"login_token_file" yaml:"login_token_file"`
	Domains        []Domain `json:"domains" yaml:"domains"`
	IPUrl          string   `json:"ip_url" yaml:"ip_url"`
	IPUrls         []string `json:"ip_urls" yaml:"ip_urls"`
	IPV6Url        string   `json:"ipv6_url" yaml:"ipv6_url"`
	IPV6Urls       []string `json:"ipv6_urls" yaml:"ipv6_urls"`
	Interval       int      `json:"interval" yaml:"interval"`
	UserAgent      string   `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`
	Socks5Proxy    string   `json:"socks5_proxy" yaml:"socks5_proxy"`
	Notify         Notify   `json:"notify" yaml:"notify"`
	Webhook        Webhook  `json:"webhook,omitempty" yaml:"webhook,omitempty"`
	IPInterface    string   `json:"ip_interface" yaml:"ip_interface"`
	IPType         string   `json:"ip_type" yaml:"ip_type"`
	Mikrotik       Mikrotik `json:"mikrotik" yaml:"mikrotik"`
	Resolver       string   `json:"resolver" yaml:"resolver"`
	UseProxy       bool     `json:"use_proxy" yaml:"use_proxy"`
	DebugInfo      bool     `json:"debug_info" yaml:"debug_info"`
	RunOnce        bool     `json:"run_once" yaml:"run_once"`
	Proxied        bool     `json:"proxied" yaml:"proxied"`
	AppKey         string   `json:"app_key" yaml:"app_key"`
	AppSecret      string   `json:"app_secret" yaml:"app_secret"`
	ConsumerKey    string   `json:"consumer_key" yaml:"consumer_key"`
	SkipSSLVerify  bool     `json:"skip_ssl_verify" yaml:"skip_ssl_verify"`
	WebPanel       WebPanel `json:"web_panel" yaml:"web_panel"`
}
