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
