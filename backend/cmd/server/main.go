package main

import (
	"flag"

	"github.com/sheet-platform/backend/internal/app"
)

func main() {
	configPath := flag.String("config", "config.yaml", "config file path")
	flag.Parse()

	app.New(*configPath).Run()
}
