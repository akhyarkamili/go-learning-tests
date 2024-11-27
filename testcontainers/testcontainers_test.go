package testcontainers

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	client "github.com/ory/hydra-client-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type HydraConfig struct {
	LoginServerHost   string
	LoginServerPort   string
	LoginServerScheme string
}

type HydraContainerPorts struct {
	AdminPort  string
	PublicPort string
}

func SetupHydra(cfg HydraConfig) (HydraContainerPorts, error) {
	urlsLogin := url.URL{
		Scheme: cfg.LoginServerScheme,
		Host:   cfg.LoginServerHost + ":" + cfg.LoginServerPort,
		Path:   "/login",
	}

	urlsConsent := url.URL{
		Scheme: cfg.LoginServerScheme,
		Host:   cfg.LoginServerHost + ":" + cfg.LoginServerPort,
		Path:   "/consent",
	}

	ctx := context.Background()
	// Define the container request with necessary environment variables and commands
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         "hydra-test-container",
			Image:        "asia-southeast2-docker.pkg.dev/sister-production/remote-dockerhub/oryd/hydra:v2.2.0",
			ExposedPorts: []string{"4444/tcp", "4445/tcp"},
			Env: map[string]string{
				"DEV":                                    "true",
				"DSN":                                    "memory",
				"LOG_LEAK_SENSITIVE_VALUES":              "false",
				"OIDC_SUBJECT_IDENTIFIERS_PAIRWISE_SALT": "abc123mustbe16chars",
				"SECRETS_SYSTEM":                         "abc123mustbe16chars",
				"URLS_SELF_PUBLIC":                       "http://127.0.0.1:4444/",
				"URLS_CONSENT":                           urlsConsent.String(),
				"URLS_LOGIN":                             urlsLogin.String(),
				"URLS_POST_LOGOUT_REDIRECT":              "https://hello",
				"SERVE_COOKIES_SAME_SITE_MODE":           "Lax",
				"SERVE_ADMIN_PORT":                       "4445",
				"SERVE_PUBLIC_PORT":                      "4444",
				"OIDC_SUBJECT_IDENTIFIERS_SUPPORTED_TYPES": "pairwise,public",
				"LOG_FORMAT":                "json",
				"LOG_REDACTION_TEXT":        "[redacted]",
				"TTL_ACCESS_TOKEN":          "1h",
				"TTL_REFRESH_TOKEN":         "24h",
				"TTL_ID_TOKEN":              "1h",
				"TTL_AUTH_CODE":             "10m",
				"TTL_LOGIN_CONSENT_REQUEST": "30m",
			},
			Cmd:        []string{"serve", "all"},
			WaitingFor: wait.ForListeningPort("4444"),
			HostConfigModifier: func(config *container.HostConfig) {
				config.AutoRemove = true
			},
		},
		Started: true,
	})

	adminPort, err := container.MappedPort(ctx, "4445")
	if err != nil {
		return HydraContainerPorts{}, err
	}
	publicPort, err := container.MappedPort(ctx, "4444")
	if err != nil {
		return HydraContainerPorts{}, err
	}

	return HydraContainerPorts{
		AdminPort:  adminPort.Port(),
		PublicPort: publicPort.Port(),
	}, nil
}

func TestHydraService(t *testing.T) {
	cfg := HydraConfig{
		LoginServerHost:   "localhost",
		LoginServerPort:   "51234",
		LoginServerScheme: "http",
	}
	ports, err := SetupHydra(cfg)
	require.NoError(t, err)

	api := createHydraAdminApiClient(ports.AdminPort, nil)

	clients, _, err := api.OAuth2API.ListOAuth2Clients(context.Background()).Execute()
	assert.NoError(t, err)
	assert.Empty(t, clients)
}

func createHydraAdminApiClient(port string, hc *http.Client) *client.APIClient {
	config := client.NewConfiguration()
	config.Servers = []client.ServerConfiguration{
		{
			URL: "http://127.0.0.1:" + port,
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
