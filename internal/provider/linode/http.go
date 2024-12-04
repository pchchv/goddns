package linode

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/pchchv/goddns/internal/settings"
	"golang.org/x/net/proxy"
	"golang.org/x/oauth2"
)

func CreateHTTPClient(conf *settings.Settings) (*http.Client, error) {
	var err error
	transport := &http.Transport{}
	if conf.UseProxy && conf.Socks5Proxy != "" {
		transport, err = applyProxy(conf.Socks5Proxy, transport)
		if err != nil {
			log.Fatalf("Error connecting to proxy: '%s'", err)
			log.Print("Continuing without proxy")
		}
	}

	if conf.LoginToken == "" {
		return nil, errors.New("LoginToken cannot be an empty string")
	}

	roundTripper := addBearerAuth(conf.LoginToken, transport)
	httpClient := http.Client{
		Timeout:   time.Second * 10,
		Transport: roundTripper,
	}

	return &httpClient, nil
}

func addBearerAuth(accessToken string, transport http.RoundTripper) http.RoundTripper {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	transportWithAuth := &oauth2.Transport{
		Source: tokenSource,
		Base:   transport,
	}

	log.Print("Using OAuth / API token to connect to DNS service")

	return transportWithAuth
}

func applyProxy(proxyAddress string, transport *http.Transport) (*http.Transport, error) {
	if proxyAddress == "" {
		log.Print("Skipping proxy: proxy address is empty string")
		return transport, nil
	}

	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		return transport, err
	}

	log.Printf("Connected to proxy : %s", proxyAddress)

	dialContext := func(_ context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}

	transport.DialContext = dialContext
	return transport, nil
}
