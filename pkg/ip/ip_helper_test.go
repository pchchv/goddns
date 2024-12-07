package ip_test

import (
	"testing"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/pkg/ip"
)

func TestGetMikrotikIP(t *testing.T) {
	t.Skip()
	conf := &settings.Settings{
		Mikrotik: settings.Mikrotik{
			Enabled:   true,
			Addr:      "http://192.168.20.1:81",
			Username:  "admin",
			Password:  "",
			Interface: "pppoe-out",
		},
	}
	helper := ip.GetIPHelperInstance(conf)
	if ip := helper.GetCurrentIP(); ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}
}

func TestGetCurrentIP(t *testing.T) {
	t.Skip()
	conf := &settings.Settings{IPUrls: []string{"https://myip.biturl.top"}}
	helper := ip.GetIPHelperInstance(conf)
	if ip := helper.GetCurrentIP(); ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}
}
