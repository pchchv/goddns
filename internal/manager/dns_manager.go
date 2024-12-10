package manager

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/pchchv/goddns/internal/handler"
	"github.com/pchchv/goddns/internal/provider"
	"github.com/pchchv/goddns/internal/server"
	"github.com/pchchv/goddns/internal/settings"
)

var managerInstance *DNSManager

type DNSManager struct {
	config      *settings.Settings
	handler     *handler.Handler
	provider    provider.IDNSProvider
	ctx         context.Context
	cancel      context.CancelFunc
	watcher     *fsnotify.Watcher
	server      *server.Server
	configPath  string
	defaultAddr string
}

func (manager *DNSManager) startServer() {
	// start the internal HTTP server
	if (manager.config.WebPanel.Addr != "" || manager.defaultAddr != ":9000") && manager.config.WebPanel.Enabled {
		manager.server = &server.Server{}
		var addr string
		if manager.config.WebPanel.Addr != "" {
			addr = manager.config.WebPanel.Addr
		} else {
			addr = manager.defaultAddr
		}
		manager.server.
			SetAddress(addr).
			SetAuthInfo(manager.config.WebPanel.Username, manager.config.WebPanel.Password).
			SetConfig(manager.config).
			SetConfigPath(manager.configPath).
			Build()

		go func() {
			if err := manager.server.Start(); err != nil {
				log.Fatalf("Failed to start the web server, error:%v", err)
			}
		}()
	} else {
		log.Println("Web panel is disabled")
	}
}

func (manager *DNSManager) initManager() error {
	log.Printf("Creating DNS handler with provider: %s", manager.config.Provider)
	dnsProvider, err := provider.GetProvider(manager.config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	manager.ctx = ctx
	manager.cancel = cancel
	manager.provider = dnsProvider
	manager.handler = &handler.Handler{}
	manager.handler.SetContext(manager.ctx)
	manager.handler.SetConfiguration(manager.config)
	manager.handler.SetProvider(manager.provider)
	manager.handler.Init()

	// if RunOnce is true, we don't need to create a file watcher and start the internal HTTP server
	if !manager.config.RunOnce {
		// create a new file watcher
		log.Println("Creating the new file watcher...")
		managerInstance.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		// monitor the configuration file changes
		managerInstance.startMonitor()
		// start the internal HTTP server
		managerInstance.startServer()
	}
	return nil
}
