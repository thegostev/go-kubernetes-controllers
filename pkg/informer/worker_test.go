package informer

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/thegostev/go-kubernetes-controllers/internal/types"
)

func TestEventWorkerProcessesEvents(t *testing.T) {
	queue := make(chan types.Event, 1)
	logger := zerolog.Nop()
	worker := NewEventWorker(queue, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker in background
	done := make(chan struct{})
	go func() {
		worker.Start(ctx)
		close(done)
	}()

	event := types.Event{
		Type:      "add",
		Namespace: "default",
		Name:      "test-deployment",
		Timestamp: time.Now(),
	}

	queue <- event
	// Give the worker a moment to process
	time.Sleep(100 * time.Millisecond)
	cancel()
	<-done
}
