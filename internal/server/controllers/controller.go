package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pchchv/goddns/internal/settings"
)

type Controller struct {
	config     *settings.Settings
	configPath string
}

func NewController(conf *settings.Settings, configPath string) *Controller {
	return &Controller{
		config:     conf,
		configPath: configPath,
	}
}

func (c *Controller) Auth(ctx fiber.Ctx) error {
	msg := "OK"
	return ctx.SendString(msg)
}
