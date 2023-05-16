package main

import "github.com/shamank/edutour_auth_service/app/internal/app"

const (
	configDir = "./configs"
)

func main() {
	app.Run(configDir)
}
