package main

import (
	"fox/config"
	"fox/internal/app"
)

func main() {
	app.Run(config.Get())
}
