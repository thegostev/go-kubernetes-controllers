package cmd

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var logLevel string

var rootCmd = &cobra.Command{
	Use:   "k8s-controller-tutorial",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your application.`,
}

func Execute() {
	level := zerolog.InfoLevel
	if parsed, err := zerolog.ParseLevel(strings.ToLower(logLevel)); err == nil {
		level = parsed
	} else {
		log.Warn().Str("level", logLevel).Msg("Invalid level, using info")
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05Z07:00"})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set log level: trace, debug, info, warn, error")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
