package main

import (
	"log"

	"github.com/Muaz717/willpower-bot/internal/pkg/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to init app: %s", err)
	}
	_ = a
}
