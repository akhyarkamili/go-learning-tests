package context_management

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGoroutineCancellation(t *testing.T) {
	t.Run("test goroutine cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		outputChan := make(chan error)
		// Act
		go sampleGoroutine(ctx, outputChan)
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		// Assert
		assert.Eventually(t, func() bool {
			select {
			case <-outputChan:
				return true
			default:
				return false
			}
		}, 3*time.Second, 300*time.Millisecond, "expected goroutine output channel to come")
	})
}
