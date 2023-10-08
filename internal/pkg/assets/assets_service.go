package assets

import "time"

type AssetsService interface {
	GetAssetById(id int64) (*Asset, error)
	GetAssetLatestPrice(asset *Asset) (*Money, error)
	GetAssetPriceAt(asset *Asset, timestamp time.Time) (*Money, error)
	StoreAssetPrice(asset *Asset, price Money, timestamp time.Time) error
}
