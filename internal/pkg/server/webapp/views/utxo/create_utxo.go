package utxo

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/currencies"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type CreateUtxoView struct {
	ledger          ledger.Ledger
	stocksService   assets.StockAssetsService
	currencyService currencies.CurrencyService
}

func NewCreateUtxoView(
	ledger ledger.Ledger,
	stocksService assets.StockAssetsService,
	currencyService currencies.CurrencyService) *CreateUtxoView {
	return &CreateUtxoView{
		ledger:          ledger,
		stocksService:   stocksService,
		currencyService: currencyService,
	}
}

func (view *CreateUtxoView) ConfigureRoutes(app *fiber.App) {
	app.Get("/utxo/create", view.Render)
	app.Get("/utxo/create/asset", view.RenderAssetControls)
}

func (view *CreateUtxoView) Render(c *fiber.Ctx) error {
	assetType := c.Query("assetType")
	assetTicker := c.Query("assetTicker")

	if assetType == "" {
		assetType = "UNKNOWN"
	}

	var assets []struct {
		Key         string
		Description string
	}

	if assetType != "UNKNOWN" {
		switch assetType {
		case "STOCK":
			stockAssets, err := view.stocksService.ListAll()
			if err != nil {
				return err
			}

			for _, asset := range stockAssets {
				assets = append(assets, struct {
					Key         string
					Description string
				}{
					Key:         asset.Ticker,
					Description: fmt.Sprintf("%s (%s)", asset.Asset.Name, asset.Ticker),
				})
			}
		}
	}

	currencies, err := view.currencyService.GetSourceCurrencies()
	if err != nil {
		return err
	}

	return c.Render("utxo/create", fiber.Map{
		"AssetType":   assetType,
		"AssetTicker": assetTicker,
		"Assets":      assets,
		"CostBasisField": fiber.Map{
			"InputName":           "costBasis",
			"InputPlaceholder":    "Cost Basis",
			"InputLabel":          "Cost Basis",
			"CurrencyName":        "costBasisCurrency",
			"CurrencyPlaceholder": "Cost Basis Currency",
			"CurrencyLabel":       "Cost Basis Currency",
			"CurrencyOptions":     currencies,
		},
		"FeesField": fiber.Map{
			"InputName":           "fees",
			"InputPlaceholder":    "Fees",
			"InputLabel":          "Fees",
			"CurrencyName":        "feesCurrency",
			"CurrencyPlaceholder": "Fees Currency",
			"CurrencyLabel":       "Fees Currency",
			"CurrencyOptions":     currencies,
		},
	})
}

func (view *CreateUtxoView) RenderAssetControls(c *fiber.Ctx) error {
	assetType := c.Query("assetType")
	if assetType == "" {
		assetType = "UNKNOWN"
	}

	c.Response().Header.Set("HX-Push-Url", fmt.Sprintf("/utxo/create?assetType=%s", assetType))

	return c.Render("utxo/create/asset", fiber.Map{
		"AssetType": assetType,
	}, "")
}

type CreateUtxoFormData struct {
	AssetType   string `form:"assetType"`
	AssetTicker string `form:"assetTicker"`
	Quantity    string `form:"quantity"`
}

func (view *CreateUtxoView) RenderPost(c *fiber.Ctx) error {
	var formData CreateUtxoFormData
	err := c.BodyParser(&formData)
	if err != nil {
		return err
	}

	if formData.AssetType == "" {
		return fmt.Errorf("asset type is required")
	}

	if formData.AssetTicker == "" {
		return fmt.Errorf("asset ticker is required")
	}

	if formData.Quantity == "" {
		return fmt.Errorf("quantity is required")
	}

	//htmxRequest := c.Request().Header.Peek("Hx-Request") != nil

	//ledger.CreateUnspentOutput()

	return nil
}
