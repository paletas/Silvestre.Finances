package views

import (
	"github.com/gofiber/fiber/v2"
	
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

func RenderIndex(c *fiber.Ctx, ledger *ledger.Ledger) error {
	return c.Render("index", fiber.Map{
		"TotalWealth": 0,
	}, "layouts/main")
}
