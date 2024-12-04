package linode

import (
	"context"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/proxy"
	"golang.org/x/oauth2"
)

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
