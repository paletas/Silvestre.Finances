package exchanges

import "github.com/paletas/silvestre.finances/internal/pkg/assets"

type ExchangeService interface {
	ListAll() ([]Exchange, error)
	ListAllByAssetType(assetType assets.AssetType) ([]Exchange, error)
}
