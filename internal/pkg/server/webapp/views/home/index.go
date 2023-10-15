package home

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
	utxo_partials "github.com/paletas/silvestre.finances/internal/pkg/server/webapp/views/utxo/partials"
)

type IndexView struct {
	ledger        ledger.Ledger
	assetsService assets.AssetsService
}

func NewIndexView(ledger ledger.Ledger, assetsService assets.AssetsService) *IndexView {
	return &IndexView{
		ledger:        ledger,
		assetsService: assetsService,
	}
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

	unspentOutputs, error := view.getUnspentOutputs(assetType)
	if error != nil {
		return error
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

	unspentOutputs, error := view.getUnspentOutputs(assetType)
	if error != nil {
		return error
	}

	c.Response().Header.Set("HX-Push-Url", fmt.Sprintf("/?assetType=%s", assetType))
	return c.Render("home/partials/index_utxolist", fiber.Map{
		"Filter": assetType,
		"Items":  unspentOutputs,
	}, "")
}

func (view *IndexView) getUnspentOutputs(assetType string) ([]utxo_partials.UnspentOutputViewModel, error) {
	var unspentOutputs []ledger.UnspentOutput
	var err error
	if assetType == "" || assetType == "ALL" {
		unspentOutputs, err = view.ledger.GetUnspentOutputs()
		if err != nil {
			return nil, err
		}
	} else {
		unspentOutputs, err = view.ledger.GetUnspentOutputsByAssetType(assets.AssetType(assetType))
		if err != nil {
			return nil, err
		}
	}

	var unspentOutputsViewModel []utxo_partials.UnspentOutputViewModel
	for _, unspentOutput := range unspentOutputs {
		assetPrice, err := view.assetsService.GetAssetLatestPriceById(unspentOutput.AssetId)
		if err != nil {
			return nil, err
		}

		unspentOutputsViewModel = append(unspentOutputsViewModel, utxo_partials.UnspentOutputViewModel{
			Id:            unspentOutput.Id,
			TransactionId: unspentOutput.TransactionId,
			Exchange:      unspentOutput.Exchange,
			Date:          unspentOutput.Date.Format("2006-01-02"),
			AssetId:       unspentOutput.AssetId,
			AssetName:     assetPrice.Asset.Name,
			Value:         assetPrice.LatestPrice.Amount * unspentOutput.Amount,
		})
	}

	return unspentOutputsViewModel, nil
}
