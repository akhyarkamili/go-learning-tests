package hydra_client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setupHydraClient(t *testing.T) {
	t.Run("test creating hydra client", func(t *testing.T) {
		// Act
		setupHydraClient("client123")
		t.Cleanup(func() {
			destroyHydraClient("client123")
		})

		// Assert
		apiClient := createHydraApiClient()
		res, _, _ := apiClient.AdminApi.ListOAuth2Clients(context.Background()).Execute()

		assert.Equal(t, len(res), 1)
		assert.Equal(t, "client123", *res[0].ClientId)
	})
}
