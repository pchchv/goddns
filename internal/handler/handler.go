package handler

import (
	"context"

	"github.com/pchchv/goddns/internal/provider"
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/pkg/ip"
	"github.com/pchchv/goddns/pkg/notification"
)

type Handler struct {
	ctx                 context.Context
	Configuration       *settings.Settings
	dnsProvider         provider.IDNSProvider
	notificationManager notification.INotificationManager
	ipManager           *ip.IPHelper
	cachedIP            string
}

func (handler *Handler) Init() {
	handler.ipManager.UpdateConfiguration(handler.Configuration)
}

func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
	handler.notificationManager = notification.GetNotificationManager(handler.Configuration)
	handler.ipManager = ip.GetIPHelperInstance(handler.Configuration)
}

func (handler *Handler) SetProvider(provider provider.IDNSProvider) {
	handler.dnsProvider = provider
}

func (handler *Handler) SetContext(ctx context.Context) {
	handler.ctx = ctx
}
