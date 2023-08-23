package main

import (
	"github.com/paletas/silvestre.finances/internal/pkg/server/webapp"
)

func main() {
	app := webapp.LaunchServer()

	defer func() {
		app.Shutdown()
	}()
}
