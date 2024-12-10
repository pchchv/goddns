package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pchchv/goddns/internal/server/controllers"
	"github.com/pchchv/goddns/internal/settings"
)

type Server struct {
	addr       string
	username   string
	password   string
	app        *fiber.App
	controller *controllers.Controller
	config     *settings.Settings
	configPath string
}
