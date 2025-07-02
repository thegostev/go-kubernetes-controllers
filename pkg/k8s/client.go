package k8s

import (
	"context"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/thegostev/go-kubernetes-controllers/internal/types"
	"github.com/thegostev/go-kubernetes-controllers/pkg/errors"
)

// Client represents a Kubernetes client
type Client struct {
	clientset *kubernetes.Clientset
	logger    zerolog.Logger
}

// NewClient creates a new Kubernetes client
func NewClient(config *types.ClientConfig) (*Client, error) {
	logger := log.With().Str("component", "k8s-client").Logger()

	// Validate configuration
	if err := config.Validate(); err != nil {
		logger.Error().Err(err).Msg("invalid client configuration")
		return nil, errors.NewConfigError("invalid client configuration", err)
	}

	// Set defaults
	config.SetDefaults()

	// Determine kubeconfig path
	kubeconfigPath := config.KubeconfigPath
	if kubeconfigPath == "" {
		kubeconfigPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	logger.Debug().Str("kubeconfig", kubeconfigPath).Msg("loading kubeconfig")

	// Load kubeconfig
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		logger.Error().Err(err).Str("kubeconfig", kubeconfigPath).Msg("failed to load kubeconfig")
		return nil, errors.NewConfigError("failed to load kubeconfig", err)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create kubernetes clientset")
		return nil, errors.NewConfigError("failed to create kubernetes clientset", err)
	}

	logger.Info().Msg("kubernetes client initialized successfully")

	return &Client{
		clientset: clientset,
		logger:    logger,
	}, nil
}

// GetClientset returns the underlying kubernetes clientset
func (c *Client) GetClientset() *kubernetes.Clientset {
	return c.clientset
}

// HealthCheck performs a basic health check
func (c *Client) HealthCheck(ctx context.Context) error {
	c.logger.Debug().Msg("performing health check")

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Get server version as health check
	_, err := c.clientset.Discovery().ServerVersion()
	if err != nil {
		c.logger.Error().Err(err).Msg("health check failed")
		return errors.NewConnectionError("health check failed", err)
	}

	c.logger.Debug().Msg("health check successful")
	return nil
}

func NewConfigOrDie() *rest.Config {
	kubeconfig := "" // TODO: optionally load from env or flag
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	return config
}
