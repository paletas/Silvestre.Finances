package main

import (
	"flag"
	"os"

	"github.com/paletas/silvestre.finances/internal/pkg/feeds/polygon"
	"github.com/paletas/silvestre.finances/internal/pkg/infrastructure/inmemory"
	"github.com/paletas/silvestre.finances/internal/pkg/server/webapp"
)

func main() {
	inmem_arg := flag.Bool("inmemory", false, "Use in-memory ledger")
	polygonApiKey := os.Getenv("POLYGON_API_KEY")

	flag.Parse()

	var app *webapp.WebApp
	if *inmem_arg {
		ledger := inmemory.CreateMemoryLedger()
		stocksService := inmemory.NewInMemoryStockAssets()
		cryptoService := inmemory.NewInMemoryCryptoAssets()
		polygonService := polygon.NewPolygonService(polygonApiKey)
		exchangeService := inmemory.NewInMemoryExchangeService()

		app = webapp.LaunchServer(&ledger, stocksService, cryptoService, polygonService, exchangeService)
	} else {
		panic("Not implemented")
	}

	defer func() {
		app.Shutdown()
	}()
}
