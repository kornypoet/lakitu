package main

import (
	_ "embed"
	"github.com/kornypoet/lakitu/api"
	"github.com/kornypoet/lakitu/cmd"
)

//go:embed VERSION
var version string

func main() {
	api.Version = version
	cmd.Execute()
}
