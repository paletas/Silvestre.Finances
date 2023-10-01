package inmemory

import "github.com/paletas/silvestre.finances/internal/pkg/assets"

type InMemoryCryptoAssets struct {
	assets []assets.CryptoAsset
}

func NewInMemoryCryptoAssets() *InMemoryCryptoAssets {
	return &InMemoryCryptoAssets{
		assets: []assets.CryptoAsset{},
	}
}

func (a *InMemoryCryptoAssets) GetKnownAssets() []assets.CryptoAsset {
	return a.assets
}

func (a *InMemoryCryptoAssets) CreateAsset(asset *assets.CryptoAsset) error {
	a.assets = append(a.assets, *asset)
	return nil
}

func (a *InMemoryCryptoAssets) GetAssetByTicker(ticker string) (*assets.CryptoAsset, error) {
	for _, asset := range a.assets {
		if asset.Ticker == ticker {
			return &asset, nil
		}
	}
	return nil, nil
}

func (a *InMemoryCryptoAssets) ListAll() ([]*assets.CryptoAsset, error) {
	var assets []*assets.CryptoAsset
	for _, asset := range a.assets {
		assets = append(assets, &asset)
	}
	return assets, nil
}
