package main

import (
	"flag"
	"log"

	"github.com/ilovepitsa/jwt_tokens/internal/app"
)

var (
	configPath = flag.String("c", "./config/config.yaml", "path to config")
)

func main() {
	if err := app.Run(*configPath); err != nil {
		log.Fatal(err)
	}
}
