package home

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type HomeController struct {
	homeView *IndexView
}

func NewHomeController(ledger ledger.Ledger, assetsService assets.AssetsService) *HomeController {
	return &HomeController{
		homeView: NewIndexView(ledger, assetsService),
	}
}

func (controller *HomeController) ConfigureRoutes(app *fiber.App) {
	controller.homeView.ConfigureRoutes(app)
}
