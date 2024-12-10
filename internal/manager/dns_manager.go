package manager

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

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

func (manager *DNSManager) Run() {
	if len(manager.config.Domains) == 0 {
		log.Println("No domain is configured, please check your configuration file")
		return
	}

	for _, domain := range manager.config.Domains {
		if manager.config.RunOnce {
			if err := manager.handler.UpdateIP(&domain); err != nil {
				log.Fatal("Error during execution:", err)
				os.Exit(1)
			}
		} else {
			// pass the context to the goroutine
			go manager.handler.LoopUpdateIP(manager.ctx, &domain)
		}
	}

	if manager.config.RunOnce {
		os.Exit(0)
	}
}

func (manager *DNSManager) Stop() {
	manager.cancel()
	// close the file watcher
	if manager.watcher != nil {
		manager.watcher.Close()
	}

	// stop the internal HTTP server
	if manager.server != nil {
		manager.server.Stop()
	}
}

func (manager *DNSManager) Restart() {
	log.Println("Restarting DNS manager...")
	manager.Stop()

	// wait for the goroutines to exit
	time.Sleep(200 * time.Millisecond)

	// re-init the manager
	if err := manager.initManager(); err != nil {
		log.Fatalf("Error during DNS manager restarting: %s", err)
	}

	manager.Run()
	log.Println("DNS manager restarted successfully")
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

func getFileName(configPath string) string {
	// get the file name from the path
	// e.g. /etc/goddns/config.json -> config.json
	return filepath.Base(configPath)
}
