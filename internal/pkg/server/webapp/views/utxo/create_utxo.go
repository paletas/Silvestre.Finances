package utxo

import (
	"fmt"
	"time"

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
	assetId := c.QueryInt("assetId")

	if assetType == "" {
		assetType = "UNKNOWN"
	}

	var existingAssets []struct {
		ID          int64
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
					ID          int64
					Description string
				}{
					ID:          asset.Id,
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
		"AssetType": assetType,
		"AssetId":   assetId,
		"Assets":    existingAssets,
		"Exchanges": exchanges,
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
	Action            string  `form:"action"`
	AssetType         string  `form:"assetType"`
	AssetId           int64   `form:"assetId"`
	Quantity          float64 `form:"quantity"`
	Exchange          string  `form:"exchange"`
	TransactionId     string  `form:"transactionId"`
	CostBasis         float64 `form:"costBasis"`
	CostBasisCurrency string  `form:"costBasisCurrency"`
	Fees              float64 `form:"fees"`
	FeesCurrency      string  `form:"feesCurrency"`
}

func (view *CreateUtxoView) RenderPost(c *fiber.Ctx) error {
	var formData CreateUtxoFormData
	err := c.BodyParser(&formData)
	if err != nil {
		return err
	}

	var existingAssets []struct {
		ID          int64
		Description string
	}

	if formData.Action == "create" {
		if formData.AssetType == "" {
			return fmt.Errorf("asset type is required")
		}
		if formData.Exchange == "" {
			return fmt.Errorf("exchange is required")
		}

		err := view.ledger.AddUnspentOutput(
			formData.TransactionId,
			formData.Exchange,
			time.Now(),
			assets.AssetType(formData.AssetType),
			formData.AssetId,
			formData.Quantity,
			assets.Money{Amount: formData.CostBasis, Currency: formData.CostBasisCurrency},
			assets.Money{Amount: formData.Fees, Currency: formData.FeesCurrency})

		if err != nil {
			return err
		}

		return c.Redirect("/")
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
						ID          int64
						Description string
					}{
						ID:          asset.Id,
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
			"AssetType": formData.AssetType,
			"AssetId":   formData.AssetId,
			"Assets":    existingAssets,
			"Exchanges": exchanges,
			"Quantity":  formData.Quantity,
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
}
