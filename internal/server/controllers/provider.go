package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pchchv/goddns/internal/utils"
)

type Provider struct {
	Provider    string `json:"provider" yaml:"provider"`
	Email       string `json:"email" yaml:"email"`
	Password    string `json:"password" yaml:"password"`
	LoginToken  string `json:"login_token" yaml:"login_token"`
	AppKey      string `json:"app_key" yaml:"app_key"`
	AppSecret   string `json:"app_secret" yaml:"app_secret"`
	ConsumerKey string `json:"consumer_key" yaml:"consumer_key"`
}

func (c *Controller) GetProvider(ctx fiber.Ctx) error {
	provider := Provider{
		Provider:    c.config.Provider,
		Email:       c.config.Email,
		Password:    c.config.Password,
		LoginToken:  c.config.LoginToken,
		AppKey:      c.config.AppKey,
		AppSecret:   c.config.AppSecret,
		ConsumerKey: c.config.ConsumerKey,
	}
	return ctx.JSON(provider)
}

func (c *Controller) GetProviderSettings(ctx fiber.Ctx) error {
	return ctx.JSON(utils.Providers)
}
