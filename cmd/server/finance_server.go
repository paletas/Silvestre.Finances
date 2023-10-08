package main

import (
	"flag"
	"os"

	"github.com/paletas/silvestre.finances/internal/pkg/feeds/polygon"
	"github.com/paletas/silvestre.finances/internal/pkg/infrastructure/inmemory"
	"github.com/paletas/silvestre.finances/internal/pkg/infrastructure/sqlite"
	assets_db "github.com/paletas/silvestre.finances/internal/pkg/infrastructure/sqlite/assets"
	ledger_db "github.com/paletas/silvestre.finances/internal/pkg/infrastructure/sqlite/ledger"
	"github.com/paletas/silvestre.finances/internal/pkg/server/webapp"
)

func main() {
	polygonApiKey := os.Getenv("POLYGON_API_KEY")

	flag.Parse()

	polygonService := polygon.NewPolygonService(polygonApiKey)
	var app *webapp.WebApp
	ledgerdb, err := ledger_db.NewLedgerDb("bin/ledger.sqlite?cache=shared&mode=memory")
	if err != nil {
		panic(err)
	}

	assetsdb, err := assets_db.NewAssetsDb("bin/assets.sqlite?cache=shared&mode=memory")
	if err != nil {
		panic(err)
	}

	defer func() {
		ledgerdb.Disconnect()
		assetsdb.Disconnect()
	}()

	assetsService := sqlite.NewDatabaseAssetsService(assetsdb)
	stocksService := sqlite.NewDatabaseStocksService(assetsdb)
	cryptoService := sqlite.NewDatabaseCryptoService(assetsdb)
	exchangeService := inmemory.NewInMemoryExchangeService()
	ledger := sqlite.NewDatabaseLedgerService(ledgerdb, assetsService)

	app = webapp.LaunchServer(ledger, stocksService, cryptoService, polygonService, exchangeService)

	defer func() {
		app.Shutdown()
	}()
}
