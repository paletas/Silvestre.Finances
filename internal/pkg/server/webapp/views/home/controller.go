package home

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type HomeController struct {
	homeView *IndexView
}

func NewHomeController(ledger ledger.Ledger) *HomeController {
	return &HomeController{
		homeView: NewIndexView(ledger),
	}
}

func (controller *HomeController) ConfigureRoutes(app *fiber.App) {
	controller.homeView.ConfigureRoutes(app)
}
