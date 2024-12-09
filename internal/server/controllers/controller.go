package controllers

import "github.com/pchchv/goddns/internal/settings"

type Controller struct {
	config     *settings.Settings
	configPath string
}
