package main

import (
	"authx-go/config"
	"authx-go/internal/app"
)

func main() {
	// Init config
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	// Run application
	if err := app.Run(cfg); err != nil {
		panic(err)
	}
}
