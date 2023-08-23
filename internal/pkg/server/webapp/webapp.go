package webapp

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
	"github.com/paletas/silvestre.finances/internal/pkg/server/webapp/views"
)

var (
	//go:embed "views/**/*.gohtml"
	viewTemplates embed.FS
)

func LaunchServer(ledger *ledger.Ledger) *fiber.App {
	app := fiber.New(fiber.Config{
		Views: html.NewFileSystem(http.FS(viewTemplates), ".gohtml"),
	})

	configureRoutes(app, ledger*ledger.Ledger)

	app.Listen(":3000")
	return app
}

func configureRoutes(app *fiber.App, ledger *ledger.Ledger) {
	app.Get("/", views.RenderIndex)
}
