package webhook

import (
	"testing"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

func TestBuildReqBody(t *testing.T) {
	t.Skip()
	w := GetWebhook(&settings.Settings{
		Webhook: settings.Webhook{
			Enabled:     true,
			URL:         "http://localhost:5000/api/v1/send",
			RequestBody: "{ \"domain\": \"{{.Domain}}\", \"ip\": \"{{.CurrentIP}}\", \"ip_type\": \"{{.IPType}}\" }",
		}})
	ret, err := w.buildReqBody("example.com", "192.168.1.1", utils.IPV4)
	if err != nil {
		t.Error(err)
	}

	if expected := `{ "domain": "example.com", "ip": "192.168.1.1", "ip_type": "IPV4" }`; ret != expected {
		t.Errorf("expected %s, got %s", expected, ret)
	}

	t.Log(ret)
}

func TestBuildReqURL(t *testing.T) {
	w := GetWebhook(&settings.Settings{
		Webhook: settings.Webhook{
			Enabled: true,
			URL:     "http://localhost:5000/api/v1/send?domain={{.Domain}}&ip={{.CurrentIP}}&ip_type={{.IPType}}",
		}})
	ret, err := w.buildReqURL("example.com", "192.168.1.1", utils.IPV4)
	if err != nil {
		t.Error(err)
	}

	if expected := "http://localhost:5000/api/v1/send?domain=example.com&ip=192.168.1.1&ip_type=IPV4"; ret != expected {
		t.Errorf("expected %s, got %s", expected, ret)
	}

	t.Log(ret)
}
