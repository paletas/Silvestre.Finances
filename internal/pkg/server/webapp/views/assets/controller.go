package assets

import (
	"github.com/gofiber/fiber/v2"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/feeds"
)

type AssetsController struct {
	stocksService assets.StockAssetsService
	cryptoService assets.CryptoAssetsService

	createStockAsset  *CreateAssetStockView
	createCryptoAsset *CreateAssetCryptoView
}

func NewAssetsController(
	stocksService assets.StockAssetsService,
	cryptoService assets.CryptoAssetsService,
	polygonFeedService feeds.StocksFeedService) *AssetsController {
	return &AssetsController{
		stocksService: stocksService,
		cryptoService: cryptoService,

		createStockAsset:  NewCreateAssetStockView(polygonFeedService, stocksService),
		createCryptoAsset: NewCreateAssetCryptoView(),
	}
}

func (controller *AssetsController) ConfigureRoutes(app *fiber.App) {
	controller.createStockAsset.ConfigureRoutes(app)
	controller.createCryptoAsset.ConfigureRoutes(app)
}
