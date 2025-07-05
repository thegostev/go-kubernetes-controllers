package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/thegostev/go-kubernetes-controllers/api/v1alpha1"
	"github.com/thegostev/go-kubernetes-controllers/internal/types"
	"github.com/thegostev/go-kubernetes-controllers/pkg/k8s"
)

var (
	frontendPageNamespace string
)

var frontendPageCmd = &cobra.Command{
	Use:   "frontendpage",
	Short: "Manage FrontendPage resources",
	Long:  `List and manage FrontendPage custom resources`,
}

var listFrontendPageCmd = &cobra.Command{
	Use:   "list",
	Short: "List FrontendPage resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listFrontendPages()
	},
}

func listFrontendPages() error {
	logger := log.With().Str("component", "frontendpage-list").Logger()

	// Create client configuration (following existing pattern)
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

	// Create list options (following existing pattern)
	listOptions := &types.ListOptions{
		Namespace: frontendPageNamespace,
		Timeout:   timeout,
	}

	// Create context
	ctx := context.Background()

	// List frontend pages
	frontendPages, err := client.ListFrontendPages(ctx, listOptions)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list frontend pages")
		return fmt.Errorf("failed to list frontend pages: %w", err)
	}

	// Display results (following existing pattern)
	if err := displayFrontendPages(frontendPages, frontendPageNamespace); err != nil {
		logger.Error().Err(err).Msg("failed to display frontend pages")
		return fmt.Errorf("failed to display frontend pages: %w", err)
	}

	return nil
}

func displayFrontendPages(frontendPages *v1alpha1.FrontendPageList, namespace string) error {
	// Simple display (following existing pattern)
	fmt.Printf("FrontendPages in namespace '%s':\n", frontendPageNamespace)
	fmt.Printf("Found %d FrontendPage(s)\n\n", len(frontendPages.Items))

	for _, page := range frontendPages.Items {
		fmt.Printf("Name: %s\n", page.Name)
		fmt.Printf("  Namespace: %s\n", page.Namespace)
		fmt.Printf("  Title: %s\n", page.Spec.Title)
		fmt.Printf("  Template: %s\n", page.Spec.Template)
		fmt.Printf("  Components: %d\n", len(page.Spec.Components))
		fmt.Printf("  Phase: %s\n", page.Status.Phase)
		if page.Status.URL != "" {
			fmt.Printf("  URL: %s\n", page.Status.URL)
		}
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(frontendPageCmd)
	frontendPageCmd.AddCommand(listFrontendPageCmd)

	listFrontendPageCmd.Flags().StringVar(&frontendPageNamespace, "namespace", "", "Namespace to list frontend pages from (default: default)")
}
