package main

import (
	"flag"
	"os"

	"github.com/paletas/silvestre.finances/internal/pkg/feeds/polygon"
	"github.com/paletas/silvestre.finances/internal/pkg/infrastructure/inmemory"
	"github.com/paletas/silvestre.finances/internal/pkg/infrastructure/sqlite"
	"github.com/paletas/silvestre.finances/internal/pkg/server/webapp"
)

func main() {
	polygonApiKey := os.Getenv("POLYGON_API_KEY")

	flag.Parse()

	polygonService := polygon.NewPolygonService(polygonApiKey)
	var app *webapp.WebApp
	db, err := sqlite.NewFinancesDb("bin/finances_db.sqlite?cache=shared&mode=memory")
	if err != nil {
		panic(err)
	}

	defer func() {
		db.Disconnect()
	}()

	exchangeService := inmemory.NewInMemoryExchangeService()
	assetsService := sqlite.NewDatabaseAssetsService(db)
	stocksService := sqlite.NewDatabaseStocksService(db)
	cryptoService := sqlite.NewDatabaseCryptoService(db)
	ledger := sqlite.NewDatabaseLedgerService(db, assetsService)

	app = webapp.LaunchServer(ledger, assetsService, stocksService, cryptoService, polygonService, exchangeService)

	defer func() {
		app.Shutdown()
	}()
}
