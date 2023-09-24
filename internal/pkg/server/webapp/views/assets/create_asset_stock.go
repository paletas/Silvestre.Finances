package assets

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/feeds"
)

type CreateAssetStockView struct {
	feedService   feeds.StocksFeedService
	assetsService assets.StockAssetsService
}

func NewCreateAssetStockView(
	feedService feeds.StocksFeedService,
	assetsService assets.StockAssetsService) *CreateAssetStockView {
	return &CreateAssetStockView{
		feedService:   feedService,
		assetsService: assetsService,
	}
}

func (view *CreateAssetStockView) ConfigureRoutes(app *fiber.App) {
	app.Get("/assets/create/stock", view.Render)
	app.Post("/assets/create/stock", view.RenderPost)
}

func (view *CreateAssetStockView) Render(c *fiber.Ctx) error {
	returnUrl := c.Query("returnUrl")

	return c.Render("assets/create/stock", fiber.Map{
		"ReturnUrl": returnUrl,
	})
}

type CreateAssetStockFormData struct {
	ReturnUrl      string `form:"returnUrl"`
	Action         string `form:"action"`
	SearchedTicker string `form:"tickerSearch"`
	StockTicker    string `form:"stockTicker"`
	StockName      string `form:"stockName"`
}

func (view *CreateAssetStockView) RenderPost(c *fiber.Ctx) error {
	var formData CreateAssetStockFormData
	err := c.BodyParser(&formData)
	if err != nil {
		return err
	}

	htmxRequest := c.Request().Header.Peek("Hx-Request") != nil

	var foundAssets []*assets.StockAsset
	if formData.Action == "" || formData.Action == "search" {
		if formData.SearchedTicker != "" && formData.StockTicker == "" {
			foundAssets, err = view.feedService.SearchStock(formData.SearchedTicker)
			if err != nil {
				return err
			}
		}
	} else {
		if formData.Action == "create" {
			if formData.StockTicker == "" {
				return errors.New("ticker is required")
			}

			newAsset := &assets.StockAsset{
				Asset: assets.Asset{
					Type: assets.StockAssetType,
					Name: formData.StockName,
				},
				Ticker: formData.StockTicker,
			}

			err = view.assetsService.CreateAsset(newAsset)
			if err != nil {
				return err
			}

			if formData.ReturnUrl == "" {
				formData.ReturnUrl = "/assets"
			} else {
				formData.ReturnUrl = fmt.Sprintf("%s&assetTicker=%s", formData.ReturnUrl, newAsset.Ticker)
			}

			if htmxRequest {
				c.Response().Header.Add("HX-Redirect", formData.ReturnUrl)
				return c.SendStatus(fiber.StatusCreated)
			} else {
				return c.Redirect(formData.ReturnUrl)
			}
		}
	}

	renderParameters := fiber.Map{
		"ReturnUrl":      formData.ReturnUrl,
		"SearchedTicker": formData.SearchedTicker,
		"Ticker":         formData.StockTicker,
		"TickerName":     formData.StockName,
		"AssetsFound":    foundAssets,
	}

	if htmxRequest {
		return c.Render("assets/create/stock", renderParameters, "")
	} else {
		return c.Render("assets/create/stock", renderParameters)
	}
}
