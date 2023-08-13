package main

import (
	"flag"
	"log"

	"github.com/juicyluv/structure-experiments/internal/app"
	"github.com/juicyluv/structure-experiments/internal/config"
)

func main() {
	configPath := flag.String("config-path", "./config", "a path for config file")
	flag.Parse()

	config.Read(*configPath)

	apl, err := app.New()
	if err != nil {
		log.Fatalf("Failed to intiialize application: %v", err)
	}

	_ = apl
}
