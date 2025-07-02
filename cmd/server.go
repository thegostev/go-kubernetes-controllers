package cmd

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"

	"github.com/thegostev/go-kubernetes-controllers/pkg/controller"
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

	// Start controller-runtime manager in background
	go func() {
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			logger.Fatal().Err(err).Msg("controller-runtime manager failed")
		}
	}()

	// --- Signal handling for graceful shutdown ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logger.Info().Str("signal", sig.String()).Msg("received shutdown signal")
	}()

	// --- FastHTTP handler with simple router ---
	handler := func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		method := string(ctx.Method())
		logger.Debug().Str("method", method).Str("path", path).Msg("Request")

		// Default root handler
		ctx.SetContentType("text/plain; charset=utf-8")
		_, _ = ctx.WriteString("Hello from FastHTTP!")
	}

	logger.Info().Str("address", addr).Msg("Server is ready to accept connections")
	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}
