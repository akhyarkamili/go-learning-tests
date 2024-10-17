package oauth_exchanger

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Exchanger struct {
	clientID     string
	clientSecret string
	tokenUrl     *url.URL
	redirectURI  string
}

func NewExchanger(clientID, clientSecret, tokenUrl, redirectURI string) (Exchanger, error) {
	u, err := url.Parse(tokenUrl)
	if err != nil {
		return Exchanger{}, err
	}
	return Exchanger{
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenUrl:     u,
		redirectURI:  redirectURI,
	}, nil
}

// ExchangeToken calls token endpoint in OIDC
func (e *Exchanger) ExchangeToken(code string) (*http.Response, error) {
	hc := http.Client{
		Timeout: 10 * time.Second,
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(e.clientID + ":" + e.clientSecret))
	req := http.Request{
		Method: "POST",
		URL:    e.tokenUrl,
		Header: map[string][]string{
			"Content-Type":  {"application/x-www-form-urlencoded"},
			"Authorization": {"Basic " + encoded},
		},
	}

	body := strings.NewReader(
		"grant_type=authorization_code&code=" + code +
			"&redirect_uri=" + e.redirectURI)
	req.Body = io.NopCloser(body)

	resp, err := hc.Do(&req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
