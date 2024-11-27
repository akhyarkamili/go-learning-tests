package acceptor

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	t.Run("serve serves an http request", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		acceptor := NewHydraAuthFlowAcceptorServer(make(chan string))
		go Serve(ctx, "51233", acceptor)
		t.Cleanup(func() {
			cancel()
		})

		assert.Eventually(t, func() bool {
			client := http.Client{}
			res, err := client.Get("http://localhost:51233/")
			return assert.NoError(t, err) && assert.Equal(t, 404, res.StatusCode)
		}, time.Second, 300*time.Millisecond, "expected status code 404")
	})
}
