package server

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
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

func (s *Server) Start() error {
	s.initRoutes()

	log.Printf("Server is listening on port: %s", s.addr)
	return s.app.Listen(s.addr)
}

func (s *Server) Stop() error {
	return s.app.ShutdownWithTimeout(200 * time.Millisecond)
}

func (s *Server) Build() {
	config := fiber.Config{}
	s.app = fiber.New(config)
	s.controller = controllers.NewController(s.config, s.configPath)
}

func (s *Server) initRoutes() {
	// set cross domain access rules
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	}))

	// middleware to rewrite paths for HTML files
	s.app.Use(func(c fiber.Ctx) error {
		// check if the request is for the API
		if strings.Contains(c.Path(), "/api/v1") {
			return c.Next()
		} else if c.Path() == "/" {
			c.Path("index.html")
			return c.Next()
		} else if strings.Contains(c.Path(), ".") {
			// check if the request is for the other resources
			return c.Next()
		}

		// rewrite request to append `.html`
		c.Path(c.Path() + ".html")
		return c.Next()
	})

	// create routes group.
	route := s.app.Group("/api/v1")
	route.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			s.username: s.password,
		},
	}))

	// register routes
	route.Get("/auth", s.controller.Auth)

	// get basic info
	route.Get("/info", s.controller.GetBasicInfo)

	// domain related routes
	route.Get("/domains", s.controller.GetDomains)
	route.Post("/domains", s.controller.AddDomain)
	route.Delete("/domains/:name", s.controller.DeleteDomain)

	// provider related routes
	route.Get("/provider", s.controller.GetProvider)
	route.Get("/provider/settings", s.controller.GetProviderSettings)
	route.Put("/provider", s.controller.UpdateProvider)

	// network related routes
	route.Get("/network", s.controller.GetNetworkSettings)
	route.Put("/network", s.controller.UpdateNetworkSettings)

	// serve embedded files
	s.app.Use("/", static.New("", static.Config{
		FS:     os.DirFS("out"),
		Browse: true,
	}))
}
