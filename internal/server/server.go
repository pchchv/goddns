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

func (s *Server) SetConfig(config *settings.Settings) *Server {
	s.config = config
	return s
}

func (s *Server) SetConfigPath(configPath string) *Server {
	s.configPath = configPath
	return s
}

func (s *Server) SetAddress(addr string) *Server {
	s.addr = addr
	return s
}

func (s *Server) SetAuthInfo(username, password string) *Server {
	s.username = username
	s.password = password
	return s
}
