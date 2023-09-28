package utxo

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/currencies"
	"github.com/paletas/silvestre.finances/internal/pkg/exchanges"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type CreateUtxoView struct {
	ledger          ledger.Ledger
	stocksService   assets.StockAssetsService
	currencyService currencies.CurrencyService
	exchangeService exchanges.ExchangeService
}

func NewCreateUtxoView(
	ledger ledger.Ledger,
	stocksService assets.StockAssetsService,
	currencyService currencies.CurrencyService,
	exchangeService exchanges.ExchangeService) *CreateUtxoView {
	return &CreateUtxoView{
		ledger:          ledger,
		stocksService:   stocksService,
		currencyService: currencyService,
		exchangeService: exchangeService,
	}
}

func (view *CreateUtxoView) ConfigureRoutes(app *fiber.App) {
	app.Get("/utxo/create", view.Render)
	app.Post("/utxo/create", view.RenderPost)
}

func (view *CreateUtxoView) Render(c *fiber.Ctx) error {
	assetType := c.Query("assetType")
	assetTicker := c.Query("assetTicker")

	if assetType == "" {
		assetType = "UNKNOWN"
	}

	var existingAssets []struct {
		Key         string
		Description string
	}

	var err error
	var exchanges []exchanges.Exchange
	if assetType != "UNKNOWN" {
		switch assetType {
		case "STOCK":
			exchanges, err = view.exchangeService.ListAllByAssetType(assets.StockAssetType)
			if err != nil {
				return err
			}

			stockAssets, err := view.stocksService.ListAll()
			if err != nil {
				return err
			}

			for _, asset := range stockAssets {
				existingAssets = append(existingAssets, struct {
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
		"Assets":      existingAssets,
		"Exchanges":   exchanges,
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

type CreateUtxoFormData struct {
	Action      string `form:"action"`
	AssetType   string `form:"assetType"`
	AssetTicker string `form:"assetTicker"`
	Quantity    string `form:"quantity"`
	Exchange    string `form:"exchange"`
}

func (view *CreateUtxoView) RenderPost(c *fiber.Ctx) error {
	var formData CreateUtxoFormData
	err := c.BodyParser(&formData)
	if err != nil {
		return err
	}

	var existingAssets []struct {
		Key         string
		Description string
	}

	if formData.Action == "create" {
		if formData.AssetType == "" {
			return fmt.Errorf("asset type is required")
		}

		if formData.AssetTicker == "" {
			return fmt.Errorf("asset ticker is required")
		}

		if formData.Quantity == "" {
			return fmt.Errorf("quantity is required")
		}
	} else {
		var exchanges []exchanges.Exchange
		if formData.AssetType != "UNKNOWN" {
			switch formData.AssetType {
			case "STOCK":
				exchanges, err = view.exchangeService.ListAllByAssetType(assets.StockAssetType)
				if err != nil {
					return err
				}

				stockAssets, err := view.stocksService.ListAll()
				if err != nil {
					return err
				}

				for _, asset := range stockAssets {
					existingAssets = append(existingAssets, struct {
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
			"AssetType":   formData.AssetType,
			"AssetTicker": formData.AssetTicker,
			"Assets":      existingAssets,
			"Exchanges":   exchanges,
			"Quantity":    formData.Quantity,
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

	return nil
}
