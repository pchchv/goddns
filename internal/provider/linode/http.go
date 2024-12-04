package linode

import (
	"log"
	"net/http"

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
