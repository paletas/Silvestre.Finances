package webapp

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/exchanges"
	"github.com/paletas/silvestre.finances/internal/pkg/feeds/polygon"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
	assets_views "github.com/paletas/silvestre.finances/internal/pkg/server/webapp/views/assets"
	home_views "github.com/paletas/silvestre.finances/internal/pkg/server/webapp/views/home"
	utxo_views "github.com/paletas/silvestre.finances/internal/pkg/server/webapp/views/utxo"
)

var (
	//go:embed "views/*/*.gohtml" "views/*/*/*.gohtml"
	viewTemplates embed.FS
)

type WebApp struct {
	ledger ledger.Ledger
	app    *fiber.App

	homeController   *home_views.HomeController
	utxoController   *utxo_views.UTXOController
	assetsController *assets_views.AssetsController
}

func LaunchServer(
	ledger ledger.Ledger,
	assetsService assets.AssetsService,
	stocksService assets.StockAssetsService,
	cryptoService assets.CryptoAssetsService,
	polygonService *polygon.PolygonService,
	exchangeService exchanges.ExchangeService) *WebApp {
	app := fiber.New(fiber.Config{
		ViewsLayout: "layouts/main",
		Views:       html.NewFileSystem(http.FS(viewTemplates), ".gohtml"),
	})

	webapp := &WebApp{
		ledger: ledger,
		app:    app,

		homeController:   home_views.NewHomeController(ledger, assetsService),
		utxoController:   utxo_views.NewUTXOController(ledger, stocksService, polygonService, exchangeService),
		assetsController: assets_views.NewAssetsController(stocksService, cryptoService, polygonService),
	}
	webapp.configureRoutes(app)

	app.Listen(":3000")
	return webapp
}

func (webapp *WebApp) Shutdown() {
	webapp.app.Shutdown()
}

func (webapp WebApp) configureRoutes(app *fiber.App) {
	webapp.homeController.ConfigureRoutes(app)
	webapp.utxoController.ConfigureRoutes(app)
	webapp.assetsController.ConfigureRoutes(app)
}
