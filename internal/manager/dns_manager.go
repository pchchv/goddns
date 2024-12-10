package manager

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/pchchv/goddns/internal/handler"
	"github.com/pchchv/goddns/internal/provider"
	"github.com/pchchv/goddns/internal/server"
	"github.com/pchchv/goddns/internal/settings"
)

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
