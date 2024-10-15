package context_management

import "context"

func sampleGoroutine(ctx context.Context, outputChan chan<- error) {
	blockingChan := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(blockingChan)
	}()
	<-blockingChan
	outputChan <- nil
}
