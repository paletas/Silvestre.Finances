package inmemory

import "github.com/paletas/silvestre.finances/internal/pkg/assets"

type InMemoryStockAssets struct {
	assets []assets.StockAsset
}

func NewInMemoryStockAssets() *InMemoryStockAssets {
	return &InMemoryStockAssets{
		assets: []assets.StockAsset{},
	}
}

func (a *InMemoryStockAssets) GetKnownAssets() []assets.StockAsset {
	return a.assets
}

func (a *InMemoryStockAssets) CreateAsset(asset *assets.StockAsset) error {
	a.assets = append(a.assets, *asset)
	return nil
}

func (a *InMemoryStockAssets) GetAssetByTicker(ticker string) (*assets.StockAsset, error) {
	for _, asset := range a.assets {
		if asset.Ticker == ticker {
			return &asset, nil
		}
	}
	return nil, nil
}

func (a *InMemoryStockAssets) ListAll() ([]*assets.StockAsset, error) {
	var assets []*assets.StockAsset
	for _, asset := range a.assets {
		assets = append(assets, &asset)
	}
	return assets, nil
}
