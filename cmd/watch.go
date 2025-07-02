package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/thegostev/go-kubernetes-controllers/internal/types"
	"github.com/thegostev/go-kubernetes-controllers/pkg/informer"
	"github.com/thegostev/go-kubernetes-controllers/pkg/k8s"
)

var (
	watchNamespace string
	watchResync    time.Duration
	watchWorkers   int
	inCluster      bool
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch Kubernetes deployments",
	Long:  `Watch deployment events using Kubernetes informer`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return watchDeployments()
	},
}

func watchDeployments() error {
	logger := log.With().Str("component", "watch-command").Logger()

	// Create client configuration
	clientConfig := &types.ClientConfig{
		KubeconfigPath: kubeconfig,
		Timeout:        timeout,
	}

	// Initialize Kubernetes client
	client, err := k8s.NewClient(clientConfig)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create kubernetes client")
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Create informer configuration
	informerConfig := &types.InformerConfig{
		Namespace:    watchNamespace,
		ResyncPeriod: watchResync,
		Workers:      watchWorkers,
	}

	// Create informer
	inf, err := informer.NewInformer(client.GetClientset(), informerConfig)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create informer")
		return fmt.Errorf("failed to create informer: %w", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info().Str("signal", sig.String()).Msg("received shutdown signal")
		cancel()
	}()

	// Start informer
	if err := inf.Start(ctx); err != nil {
		logger.Error().Err(err).Msg("failed to start informer")
		return fmt.Errorf("failed to start informer: %w", err)
	}

	logger.Info().Msg("watching deployment events (press Ctrl+C to stop)")

	// Wait for context cancellation
	<-ctx.Done()

	// Stop informer
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer stopCancel()

	if err := inf.Stop(stopCtx); err != nil {
		logger.Error().Err(err).Msg("failed to stop informer")
		return fmt.Errorf("failed to stop informer: %w", err)
	}

	logger.Info().Msg("watch stopped successfully")
	return nil
}

func init() {
	rootCmd.AddCommand(watchCmd)

	// Add flags
	watchCmd.Flags().StringVar(&watchNamespace, "namespace", "default", "Namespace to watch")
	watchCmd.Flags().DurationVar(&watchResync, "resync", 10*time.Minute, "Resync period")
	watchCmd.Flags().IntVar(&watchWorkers, "workers", 2, "Number of event workers")
	watchCmd.Flags().BoolVar(&inCluster, "in-cluster", false, "Use in-cluster authentication")
}
