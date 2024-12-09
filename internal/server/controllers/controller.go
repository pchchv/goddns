package controllers

import "github.com/pchchv/goddns/internal/settings"

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
