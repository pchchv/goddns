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
