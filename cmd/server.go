package cmd

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
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
	log.Info().Str("component", "server").Int("port", port).Msg("Starting server")

	handler := func(ctx *fasthttp.RequestCtx) {
		log.Debug().Str("method", string(ctx.Method())).Str("path", string(ctx.Path())).Msg("Request")
		ctx.SetContentType("text/plain; charset=utf-8")
		if _, err := ctx.WriteString("Hello from FastHTTP!"); err != nil {
			log.Error().Err(err).Msg("Write failed")
		}
	}

	log.Info().Str("component", "server").Str("address", addr).Msg("Server is ready to accept connections")

	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
