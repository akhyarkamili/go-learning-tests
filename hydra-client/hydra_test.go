package hydra_client

import (
	"context"
	"learning-tests/hydra-client/acceptor"
	hydra_client "learning-tests/hydra-client/hydra-client"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_setupHydraClient(t *testing.T) {
	t.Run("test creating hydra client", func(t *testing.T) {
		// Act
		setupOIDCClient("client123", "secret123")
		t.Cleanup(func() {
			destroyOIDCClient("client123")
		})

		// Assert
		apiClient := hydra_client.CreateHydraApiClient("4445", nil)
		res, _, _ := apiClient.AdminApi.ListOAuth2Clients(context.Background()).Execute()

		assert.Equal(t, len(res), 1)
		assert.Equal(t, "client123", *res[0].ClientId)
	})
}

func Test_setupAuthentication(t *testing.T) {
	t.Run("test setup authentication", func(t *testing.T) {
		// Arrange
		id := "client123"
		secret := "secret123"
		acceptorPort := "51234"

		setupOIDCClient(id, secret)
		t.Cleanup(func() {
			destroyOIDCClient(id)
		})

		ctx, cancel := context.WithCancel(context.Background())

		authCode := make(chan string)
		server := acceptor.NewHydraAuthFlowAcceptorServer(authCode)
		go acceptor.Serve(ctx, acceptorPort, server)
		t.Cleanup(func() {
			cancel()
		})

		// Act
		err := setupAuthentication(id)
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			select {
			case code := <-authCode:
				return code != ""
			default:
				return false
			}
		}, 3*time.Second, 300*time.Millisecond, "expected auth code to be received")

	})
}
