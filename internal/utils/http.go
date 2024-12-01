package utils

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/pchchv/goddns/internal/settings"
	"golang.org/x/net/proxy"
)

// GetHTTPClient creates the HTTP client and return it.
func GetHTTPClient(conf *settings.Settings) *http.Client {
	client := &http.Client{
		Timeout: time.Second * DefaultTimeout,
	}

	if conf.UseProxy && conf.Socks5Proxy != "" {
		log.Println("use socks5 proxy:" + conf.Socks5Proxy)
		dialer, err := proxy.SOCKS5("tcp", conf.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			log.Fatal("can't connect to the proxy:", err)
			return nil
		}

		dialContext := func(_ context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		}

		httpTransport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.SkipSSLVerify},
		}
		client.Transport = httpTransport
		httpTransport.DialContext = dialContext
	} else {
		httpTransport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.SkipSSLVerify},
		}
		client.Transport = httpTransport
	}

	return client
}
