package controllers

type Provider struct {
	Provider    string `json:"provider" yaml:"provider"`
	Email       string `json:"email" yaml:"email"`
	Password    string `json:"password" yaml:"password"`
	LoginToken  string `json:"login_token" yaml:"login_token"`
	AppKey      string `json:"app_key" yaml:"app_key"`
	AppSecret   string `json:"app_secret" yaml:"app_secret"`
	ConsumerKey string `json:"consumer_key" yaml:"consumer_key"`
}
