package cmd

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

var port int

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server",
	Long: `Start a FastHTTP server with configurable port and log level.
The server will respond with "Hello from FastHTTP!" to any request.`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Add port flag
	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
}

func startServer() {
	addr := ":" + strconv.Itoa(port)

	log.Info().
		Str("component", "server").
		Int("port", port).
		Str("address", addr).
		Msg("Starting FastHTTP server")

	// Define request handler
	handler := func(ctx *fasthttp.RequestCtx) {
		log.Debug().
			Str("method", string(ctx.Method())).
			Str("path", string(ctx.Path())).
			Str("remote_addr", ctx.RemoteAddr().String()).
			Msg("Request received")

		// Set response headers
		ctx.SetContentType("text/plain; charset=utf-8")

		// Send response
		ctx.WriteString("Hello from FastHTTP!")

		log.Debug().
			Str("method", string(ctx.Method())).
			Str("path", string(ctx.Path())).
			Int("status_code", ctx.Response.StatusCode()).
			Msg("Response sent")
	}

	// Start server
	log.Info().
		Str("component", "server").
		Str("address", addr).
		Msg("Server is ready to accept connections")

	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		log.Fatal().
			Err(err).
			Str("component", "server").
			Str("address", addr).
			Msg("Failed to start server")
	}
}
