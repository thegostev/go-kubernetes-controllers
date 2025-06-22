/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-controller-tutorial",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Configure logging based on log level flag
	configureLogging()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func configureLogging() {
	// Parse log level
	level, err := zerolog.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		level = zerolog.InfoLevel
		log.Warn().Str("invalid_level", logLevel).Msg("Invalid log level, defaulting to info")
	}

	// Set global log level
	zerolog.SetGlobalLevel(level)

	// Configure pretty console output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05Z07:00"})

	log.Info().Str("level", level.String()).Msg("Logging configured")
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-controller-tutorial.yaml)")

	// Add log level flag
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set log level: trace, debug, info, warn, error")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
