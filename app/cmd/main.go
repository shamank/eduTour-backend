package main

import "github.com/shamank/eduTour-backend/app/internal/app"

const (
	configDir = "./configs"
)

func main() {
	app.Run(configDir)
}
