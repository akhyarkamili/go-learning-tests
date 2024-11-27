package context_management

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGoroutineCancellation(t *testing.T) {
	t.Run("cancel sends into ctx.Done in another goroutine", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// Act
		outputChan := make(chan error)
		go func(ctx context.Context, outputChan chan<- error) {
			<-ctx.Done()
			outputChan <- nil
		}(ctx, outputChan)
		cancel()

		// Assert
		assert.Eventually(t, func() bool {
			select {
			case <-outputChan:
				return true
			default:
				return false
			}
		}, 1*time.Second, 100*time.Millisecond, "expected goroutine output channel to come")
	})
}
