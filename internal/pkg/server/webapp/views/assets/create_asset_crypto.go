package assets

import "github.com/gofiber/fiber/v2"

type CreateAssetCryptoView struct {
}

func NewCreateAssetCryptoView() *CreateAssetCryptoView {
	return &CreateAssetCryptoView{}
}

func (view *CreateAssetCryptoView) ConfigureRoutes(app *fiber.App) {
	app.Get("/assets/create/crypto", view.Render)
}

func (view *CreateAssetCryptoView) Render(c *fiber.Ctx) error {
	return c.Render("assets/create/crypto", fiber.Map{})
}
