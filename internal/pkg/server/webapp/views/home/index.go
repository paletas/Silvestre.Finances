package home

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type IndexView struct {
	ledger ledger.Ledger
}

func NewIndexView(ledger ledger.Ledger) *IndexView {
	return &IndexView{ledger: ledger}
}

func (view *IndexView) ConfigureRoutes(app *fiber.App) {
	app.Get("/", view.Render)
	app.Get("/partial/utxo_list", view.RenderUtxoList)
}

func (view *IndexView) Render(c *fiber.Ctx) error {
	assetType := c.Query("assetType")

	totalWealth, error := view.ledger.GetTotalWealth()
	if error != nil {
		return error
	}

	var unspentOutputs []ledger.UnspentOutput
	if assetType == "" {
		assetType = "ALL"
		unspentOutputs, error = view.ledger.GetUnspentOutputs()
		if error != nil {
			return error
		}
	} else {
		c.Response().Header.Set("HX-Push-Url", fmt.Sprintf("/?assetType=%s", assetType))

		unspentOutputs, error = view.ledger.GetUnspentOutputsByAssetType(assets.AssetType(assetType))
		if error != nil {
			return error
		}
	}

	return c.Render("home/index", fiber.Map{
		"TotalWealth": totalWealth,
		"Filter":      assetType,
		"Items":       unspentOutputs,
	})
}

func (view *IndexView) RenderUtxoList(c *fiber.Ctx) error {
	assetType := c.Query("assetType")

	if assetType == "" {
		return c.SendStatus(400)
	}

	unspentOutputs, error := view.ledger.GetUnspentOutputsByAssetType(assets.AssetType(assetType))
	if error != nil {
		return error
	}

	c.Response().Header.Set("HX-Push-Url", fmt.Sprintf("/?assetType=%s", assetType))
	return c.Render("home/partials/index_utxolist", fiber.Map{
		"Filter": assetType,
		"Items":  unspentOutputs,
	}, "")
}
