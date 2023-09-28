package utxo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/currencies"
	"github.com/paletas/silvestre.finances/internal/pkg/exchanges"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type UTXOController struct {
	utxoView *CreateUtxoView
}

func NewUTXOController(
	ledger ledger.Ledger,
	stocksService assets.StockAssetsService,
	currencyService currencies.CurrencyService,
	exchangeService exchanges.ExchangeService,
) *UTXOController {
	return &UTXOController{
		utxoView: NewCreateUtxoView(ledger, stocksService, currencyService, exchangeService),
	}
}

func (controller *UTXOController) ConfigureRoutes(app *fiber.App) {
	controller.utxoView.ConfigureRoutes(app)
}
