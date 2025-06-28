package informer

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/yourusername/k8s-controller-tutorial/internal/types"
)

// EventWorker processes events asynchronously
type EventWorker struct {
	queue  chan types.Event
	logger zerolog.Logger
}

// NewEventWorker creates a new event worker
func NewEventWorker(queue chan types.Event, logger zerolog.Logger) *EventWorker {
	return &EventWorker{
		queue:  queue,
		logger: logger,
	}
}

// Start starts the event worker
func (w *EventWorker) Start(ctx context.Context) {
	w.logger.Debug().Msg("starting event worker")

	for {
		select {
		case event, ok := <-w.queue:
			if !ok {
				w.logger.Debug().Msg("event queue closed, stopping worker")
				return
			}
			w.processEvent(event)
		case <-ctx.Done():
			w.logger.Debug().Msg("context cancelled, stopping worker")
			return
		}
	}
}

// processEvent processes a single event
func (w *EventWorker) processEvent(event types.Event) {
	w.logger.Info().
		Str("type", event.Type).
		Str("namespace", event.Namespace).
		Str("name", event.Name).
		Time("timestamp", event.Timestamp).
		Msg("processing deployment event")

	// In a real implementation, you would do more processing here
	// For now, we just log the event
}
