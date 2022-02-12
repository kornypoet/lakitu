package cmd

import (
	"fmt"
	"os"

	"github.com/kornypoet/lakitu/api"
	"github.com/spf13/cobra"
)

var Bind string
var Debug bool
var Logging bool
var Port string

var rootCmd = &cobra.Command{
	Use:   "lakitu",
	Short: "Simple HTTP Server",
	Run:   run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&Bind, "bind", "b", "localhost", "Address to bind server to")
	rootCmd.Flags().BoolVarP(&Debug, "debug", "d", false, "Enable debug mode")
	rootCmd.Flags().BoolVarP(&Logging, "logging", "l", true, "Enable request logs")
	rootCmd.Flags().StringVarP(&Port, "port", "p", "8080", "Port to listen on")
}

func run(cmd *cobra.Command, args []string) {
	router := api.Router(Debug, Logging)
	address := fmt.Sprintf("%s:%s", Bind, Port)
	router.Run(address)
}
