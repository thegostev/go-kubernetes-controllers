package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/thegostev/go-kubernetes-controllers/internal/types"
	"github.com/thegostev/go-kubernetes-controllers/pkg/k8s"
)

var (
	kubeconfig string
	namespace  string
	timeout    time.Duration
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Kubernetes deployments",
	Long:  `List all deployments in the specified namespace`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listDeployments()
	},
}

func listDeployments() error {
	logger := log.With().Str("component", "list-command").Logger()

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

	// Create list options
	listOptions := &types.ListOptions{
		Namespace: namespace,
		Timeout:   timeout,
	}

	// Create context
	ctx := context.Background()

	// Perform health check
	if err := client.HealthCheck(ctx); err != nil {
		logger.Error().Err(err).Msg("health check failed")
		return fmt.Errorf("health check failed: %w", err)
	}

	// List deployments
	deployments, err := client.ListDeployments(ctx, listOptions)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list deployments")
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	// Display results
	if err := displayDeployments(deployments); err != nil {
		logger.Error().Err(err).Msg("failed to display deployments")
		return fmt.Errorf("failed to display deployments: %w", err)
	}

	return nil
}

func displayDeployments(deployments interface{}) error {
	// For now, just print a simple message
	// In a real implementation, this would format the deployments properly
	fmt.Println("Deployments listed successfully")
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file (default: ~/.kube/config)")
	listCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace to list deployments from")
	listCmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "Timeout for operations")
}
