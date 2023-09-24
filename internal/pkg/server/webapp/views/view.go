package views

import (
	"github.com/gofiber/fiber/v2"
)

type View interface {
	ConfigureRoutes(app *fiber.App) error
}
