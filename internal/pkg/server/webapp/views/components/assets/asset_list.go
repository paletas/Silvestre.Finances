package assets

import (
	"github.com/gofiber/fiber/v2"
)

func Render(c *fiber.Ctx) error {
	return c.Render("components/assets/asset_list", fiber.Map{
		"TotalWealth": 0,
	})
}
