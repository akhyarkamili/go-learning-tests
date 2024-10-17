package oauth_exchanger

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExchangeTokenWithInvalidCodeReturns400(t *testing.T) {
	// Arrange
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tokenUrl := os.Getenv("TOKEN_URL")
	redirectURI := os.Getenv("REDIRECT_URI")

	exchanger, err := NewExchanger(clientID, clientSecret, tokenUrl, redirectURI)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	code := ""

	// Act
	resp, err := exchanger.ExchangeToken(code)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	content := make(map[string]any)
	err = json.Unmarshal(body, &content)
	assert.NoError(t, err)
	assert.Contains(t, content, "error")
}
