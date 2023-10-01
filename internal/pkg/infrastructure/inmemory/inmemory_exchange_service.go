package inmemory

import (
	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/exchanges"
)

type InMemoryExchangeService struct {
	exchanges []exchanges.Exchange
}

func NewInMemoryExchangeService() *InMemoryExchangeService {
	return &InMemoryExchangeService{
		exchanges: []exchanges.Exchange{
			{
				ID:        1,
				Name:      "Binance",
				AssetType: []assets.AssetType{assets.CryptoAssetType},
			},
			{
				ID:        2,
				Name:      "Revolut",
				AssetType: []assets.AssetType{assets.StockAssetType, assets.CryptoAssetType},
			},
			{
				ID:        3,
				Name:      "Degiro",
				AssetType: []assets.AssetType{assets.StockAssetType},
			},
			{
				ID:        4,
				Name:      "Trading212",
				AssetType: []assets.AssetType{assets.StockAssetType},
			},
		},
	}
}

func (service *InMemoryExchangeService) ListAll() ([]exchanges.Exchange, error) {
	return service.exchanges, nil
}

func (service *InMemoryExchangeService) ListAllByAssetType(assetType assets.AssetType) ([]exchanges.Exchange, error) {
	var exchanges []exchanges.Exchange

	for _, exchange := range service.exchanges {
		for _, exchangeAssetType := range exchange.AssetType {
			if exchangeAssetType == assetType {
				exchanges = append(exchanges, exchange)
			}
		}
	}

	return exchanges, nil
}
