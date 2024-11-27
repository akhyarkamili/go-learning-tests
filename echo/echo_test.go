package echo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	a := ":1323"
	go Serve(a)

	// Add a test case to test the echo server
	// Use the http package to send a GET request to the echo server
	client := http.Client{}
	resp, err := client.Get("http://localhost:1323/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
