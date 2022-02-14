package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kornypoet/lakitu/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Assets  string
	Bind    string
	Debug   bool
	Logging bool
	Port    string
)

var rootCmd = &cobra.Command{
	Use:     "lakitu",
	Short:   "Simple HTTP Server",
	Run:     run,
	Version: "embedded",
}

func Execute() {
	rootCmd.Version = api.Version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&Assets, "assets", "a", "", "Directory to store assets in")
	rootCmd.Flags().StringVarP(&Bind, "bind", "b", "localhost", "Address to bind server to")
	rootCmd.Flags().BoolVarP(&Debug, "debug", "d", false, "Enable debug mode")
	rootCmd.Flags().BoolVarP(&Logging, "logging", "l", true, "Enable request logs")
	rootCmd.Flags().StringVarP(&Port, "port", "p", "8080", "Port to listen on")
}

func run(cmd *cobra.Command, args []string) {
	if Debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Info("Running lakitu")
	if Assets == "" {
		wd, _ := os.Getwd()
		Assets = filepath.Join(wd, "assets")
		log.Debug("asset dir unspecified")
	}
	log.Debugf("creating asset dir at: %s", Assets)
	os.MkdirAll(Assets, 0700)
	api.AssetDir = Assets
	router := api.Router(Logging)
	address := fmt.Sprintf("%s:%s", Bind, Port)
	log.Infof("Starting server at %s", address)
	router.Run(address)
}
