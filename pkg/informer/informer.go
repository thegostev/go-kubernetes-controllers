package informer

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/yourusername/k8s-controller-tutorial/internal/types"
	"github.com/yourusername/k8s-controller-tutorial/pkg/errors"
)

// Informer represents a Kubernetes deployment informer
type Informer struct {
	clientset  *kubernetes.Clientset
	config     *types.InformerConfig
	logger     zerolog.Logger
	indexer    cache.Indexer
	informer   cache.SharedIndexInformer
	eventQueue chan types.Event
	health     *types.InformerHealth
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	workers    []*EventWorker
	wg         sync.WaitGroup
}

// NewInformer creates a new deployment informer
func NewInformer(clientset *kubernetes.Clientset, config *types.InformerConfig) (*Informer, error) {
	logger := log.With().Str("component", "informer").Logger()

	if err := config.Validate(); err != nil {
		logger.Error().Err(err).Msg("invalid informer configuration")
		return nil, errors.NewConfigError("invalid informer configuration", err)
	}
	config.SetDefaults()

	ctx, cancel := context.WithCancel(context.Background())
	eventQueue := make(chan types.Event, config.EventBufferSize)
	health := &types.InformerHealth{
		IsHealthy: true,
		LastSync:  time.Now(),
		Workers:   config.Workers,
	}

	// Use the shared informer factory for deployments
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		config.ResyncPeriod,
		informers.WithNamespace(config.Namespace),
	)
	deploymentInformer := factory.Apps().V1().Deployments().Informer()

	informer := &Informer{
		clientset:  clientset,
		config:     config,
		logger:     logger,
		informer:   deploymentInformer,
		indexer:    deploymentInformer.GetIndexer(),
		eventQueue: eventQueue,
		health:     health,
		ctx:        ctx,
		cancel:     cancel,
	}

	// Register event handlers
	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    informer.handleAdd,
		UpdateFunc: informer.handleUpdate,
		DeleteFunc: informer.handleDelete,
	})

	// Create event workers
	informer.workers = make([]*EventWorker, config.Workers)
	for j := 0; j < config.Workers; j++ {
		informer.workers[j] = NewEventWorker(eventQueue, logger)
	}

	logger.Info().
		Str("namespace", config.Namespace).
		Dur("resyncPeriod", config.ResyncPeriod).
		Int("workers", config.Workers).
		Msg("informer initialized successfully")

	return informer, nil
}

// Start starts the informer
func (i *Informer) Start(ctx context.Context) error {
	i.logger.Info().Msg("starting informer")

	// Start event workers
	for _, worker := range i.workers {
		i.wg.Add(1)
		go func(w *EventWorker) {
			defer i.wg.Done()
			w.Start(ctx)
		}(worker)
	}

	// Start informer
	go func() {
		i.informer.Run(ctx.Done())
	}()

	// Wait for cache sync
	if !cache.WaitForCacheSync(ctx.Done(), i.informer.HasSynced) {
		return errors.NewWatchError("failed to sync cache", nil)
	}

	i.logger.Info().Msg("informer started successfully")
	return nil
}

// Stop stops the informer
func (i *Informer) Stop(ctx context.Context) error {
	i.logger.Info().Msg("stopping informer")
	i.cancel()
	close(i.eventQueue)
	done := make(chan struct{})
	go func() {
		i.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		i.logger.Info().Msg("informer stopped successfully")
		return nil
	case <-ctx.Done():
		return errors.NewWatchError("timeout waiting for informer to stop", ctx.Err())
	}
}

// Health returns the informer health status
func (i *Informer) Health() types.InformerHealth {
	i.mu.RLock()
	defer i.mu.RUnlock()
	health := *i.health
	health.CacheSize = len(i.indexer.ListKeys())
	return health
}

// handleAdd handles deployment add events
func (i *Informer) handleAdd(obj interface{}) {
	deployment := obj.(*appsv1.Deployment)
	event := types.Event{
		Type:      "add",
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
		Timestamp: time.Now(),
		Object:    deployment,
	}
	i.logger.Debug().
		Str("type", "add").
		Str("namespace", deployment.Namespace).
		Str("name", deployment.Name).
		Msg("deployment added")
	select {
	case i.eventQueue <- event:
	default:
		i.logger.Warn().Msg("event queue full, dropping add event")
	}
}

// handleUpdate handles deployment update events
func (i *Informer) handleUpdate(oldObj, newObj interface{}) {
	newDeployment := newObj.(*appsv1.Deployment)
	event := types.Event{
		Type:      "update",
		Namespace: newDeployment.Namespace,
		Name:      newDeployment.Name,
		Timestamp: time.Now(),
		Object:    newDeployment,
	}
	i.logger.Debug().
		Str("type", "update").
		Str("namespace", newDeployment.Namespace).
		Str("name", newDeployment.Name).
		Msg("deployment updated")
	select {
	case i.eventQueue <- event:
	default:
		i.logger.Warn().Msg("event queue full, dropping update event")
	}
}

// handleDelete handles deployment delete events
func (i *Informer) handleDelete(obj interface{}) {
	deployment := obj.(*appsv1.Deployment)
	event := types.Event{
		Type:      "delete",
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
		Timestamp: time.Now(),
		Object:    deployment,
	}
	i.logger.Debug().
		Str("type", "delete").
		Str("namespace", deployment.Namespace).
		Str("name", deployment.Name).
		Msg("deployment deleted")
	select {
	case i.eventQueue <- event:
	default:
		i.logger.Warn().Msg("event queue full, dropping delete event")
	}
}

// GetDeployment retrieves a deployment from cache
func (i *Informer) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	key := namespace + "/" + name
	obj, exists, err := i.indexer.GetByKey(key)
	if err != nil {
		return nil, errors.NewCacheError("failed to get deployment from cache", err)
	}
	if !exists {
		return nil, errors.NewCacheError("deployment not found in cache", nil)
	}
	return obj.(*appsv1.Deployment), nil
}

// ListDeployments lists all deployments in cache
func (i *Informer) ListDeployments() ([]*appsv1.Deployment, error) {
	objs := i.indexer.List()
	deployments := make([]*appsv1.Deployment, len(objs))
	for i, obj := range objs {
		deployments[i] = obj.(*appsv1.Deployment)
	}
	return deployments, nil
}
