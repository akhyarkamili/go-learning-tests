package hydra_client

import (
	"net/http"
	"time"

	"github.com/ory/hydra-client-go"
)

func CreateHydraApiClient(port string, hc *http.Client) *client.APIClient {
	config := client.NewConfiguration()
	config.Servers = []client.ServerConfiguration{
		{
			URL: "http://localhost:" + port,
		},
	}
	if hc != nil {
		config.HTTPClient = hc
	} else {
		config.HTTPClient = &http.Client{
			Timeout: 1 * time.Second,
		}
	}

	apiClient := client.NewAPIClient(config)
	return apiClient
}
