package cmd

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"

	"github.com/thegostev/go-kubernetes-controllers/internal/types"
	"github.com/thegostev/go-kubernetes-controllers/pkg/controller"
	"github.com/thegostev/go-kubernetes-controllers/pkg/informer"
	"github.com/thegostev/go-kubernetes-controllers/pkg/k8s"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var port int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server",
	Run:   func(cmd *cobra.Command, args []string) { startServer() },
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
}

func startServer() {
	addr := ":" + strconv.Itoa(port)
	logger := log.With().Str("component", "server").Int("port", port).Logger()

	// --- Controller-runtime manager and controller setup ---
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	mgr, err := ctrl.NewManager(k8s.NewConfigOrDie(), ctrl.Options{
		Scheme: k8s.NewScheme(),
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to start controller-runtime manager")
	}
	if err := controller.SetupDeploymentController(mgr); err != nil {
		logger.Fatal().Err(err).Msg("unable to setup deployment controller")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start controller-runtime manager in background
	go func() {
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			logger.Fatal().Err(err).Msg("controller-runtime manager failed")
		}
	}()

	// --- Informer and client setup ---
	clientConfig := &types.ClientConfig{
		KubeconfigPath: "", // default kubeconfig
		Timeout:        30 * time.Second,
	}
	client, err := k8s.NewClient(clientConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create kubernetes client")
	}
	informerConfig := &types.InformerConfig{
		Namespace:       "default",
		ResyncPeriod:    10 * time.Minute,
		Workers:         2,
		MaxCacheSize:    1000,
		MaxConnections:  10,
		EventBufferSize: 100,
	}
	inf, err := informer.NewInformer(client.GetClientset(), informerConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create informer")
	}

	// Start informer in background
	go func() {
		if err := inf.Start(ctx); err != nil {
			logger.Fatal().Err(err).Msg("informer failed")
		}
	}()

	// --- Signal handling for graceful shutdown ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logger.Info().Str("signal", sig.String()).Msg("received shutdown signal")
		cancel()
	}()

	// --- FastHTTP handler with simple router ---
	handler := func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		method := string(ctx.Method())
		logger.Debug().Str("method", method).Str("path", path).Msg("Request")

		if path == "/api/deployments" && method == fasthttp.MethodGet {
			deployments, err := inf.ListDeployments()
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				ctx.SetContentType("application/json")
				_, _ = ctx.Write([]byte(`{"error":"failed to list deployments"}`))
				return
			}
			ctx.SetContentType("application/json")
			enc := json.NewEncoder(ctx)
			enc.SetIndent("", "  ")
			if err := enc.Encode(deployments); err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				ctx.SetContentType("application/json")
				_, _ = ctx.Write([]byte(`{"error":"failed to encode deployments"}`))
			}
			return
		}

		// Default root handler
		ctx.SetContentType("text/plain; charset=utf-8")
		_, _ = ctx.WriteString("Hello from FastHTTP!")
	}

	logger.Info().Str("address", addr).Msg("Server is ready to accept connections")
	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		logger.Fatal().Err(err).Msg("Server failed")
	}

	// On shutdown, stop informer
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()
	_ = inf.Stop(stopCtx)
}
