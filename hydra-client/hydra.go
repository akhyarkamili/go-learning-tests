package hydra_client

import (
	"context"
	"fmt"
	hydra_client "learning-tests/hydra-client/hydra-client"
	"net/http"
	"net/http/cookiejar"
	url "net/url"

	client "github.com/ory/hydra-client-go"
)

func setupAuthentication(clientId string) error {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	hc := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Printf("redirecting to: %s\n", req.URL)
			fmt.Printf("Cookie jar: %v\n", jar)
			return nil
		},
		Jar: jar,
	}

	u := url.URL{Scheme: "http", Host: "127.0.0.1:4444", Path: "/oauth2/auth"}
	var v = make(url.Values)
	v.Add("response_type", "code")
	v.Add("scope", "openid")
	v.Add("state", "abcdefghijklkmn21341234")
	v.Add("client_id", clientId)
	v.Add("redirect_uri", "http://localhost:51234/callback")
	u.RawQuery = v.Encode()
	resp, err := hc.Get(u.String())
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if resp.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("failed to authenticate: %d", resp.StatusCode)
	}

}

func setupOIDCClient(clientId, clientSecret string) error {
	apiClient := hydra_client.CreateHydraApiClient("4445", nil)

	req := apiClient.AdminApi.CreateOAuth2Client(context.Background())
	oc := client.NewOAuth2ClientWithDefaults()
	oc.SetClientId(clientId)
	oc.SetClientName("client123")
	oc.SetClientSecret(clientSecret)
	oc.SetGrantTypes([]string{"authorization_code"})
	oc.SetRedirectUris([]string{"http://localhost:51234/callback"})
	oc.SetResponseTypes([]string{"code"})
	oc.SetScope("openid")
	oc.SetTokenEndpointAuthMethod("client_secret_basic")
	req = req.OAuth2Client(*oc)

	_, _, err := req.Execute()
	if err != nil {
		return fmt.Errorf("failed to create hydra client: %w", err)
	}

	return nil
}

func destroyOIDCClient(clientId string) error {
	apiClient := hydra_client.CreateHydraApiClient("4445", nil)
	req := apiClient.AdminApi.DeleteOAuth2Client(context.Background(), clientId)
	_, err := req.Execute()
	if err != nil {
		return fmt.Errorf("failed to delete hydra client: %v", err)
	}
	return nil
}
