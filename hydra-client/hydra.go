package hydra_client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	client "github.com/ory/hydra-client-go"
)

func setupAuthentication() {
	log.Fatalf("todo")
}

func setupHydraClient(clientId string) error {
	apiClient := createHydraApiClient()

	req := apiClient.AdminApi.CreateOAuth2Client(context.Background())
	oc := client.NewOAuth2ClientWithDefaults()
	oc.SetClientId(clientId)
	oc.SetClientName("client123")
	oc.SetClientSecret("secret123")
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

func destroyHydraClient(clientId string) error {
	apiClient := createHydraApiClient()
	req := apiClient.AdminApi.DeleteOAuth2Client(context.Background(), clientId)
	_, err := req.Execute()
	if err != nil {
		return fmt.Errorf("failed to delete hydra client: %v", err)
	}
	return nil
}

func createHydraApiClient() *client.APIClient {
	config := client.NewConfiguration()
	config.Servers = []client.ServerConfiguration{
		{
			URL: "http://localhost:4445",
		},
	}
	config.HTTPClient = &http.Client{
		Timeout: 1 * time.Second,
	}

	apiClient := client.NewAPIClient(config)
	return apiClient
}

func createHydraClientBody() []byte {
	body := map[string]interface{}{
		"client_id":                  "client123",
		"client_secret":              "secret123",
		"grant_types":                []string{"authorization_code"},
		"redirect_uris":              []string{"http://localhost:51234/callback"},
		"response_types":             []string{"code"},
		"scope":                      "openid",
		"token_endpoint_auth_method": "client_secret_basic",
	}
	b, _ := json.Marshal(body)
	return b
}
