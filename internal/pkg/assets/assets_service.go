package assets

import "time"

type AssetsService interface {
	GetAssetById(id int64) (*Asset, error)
	GetAssetLatestPriceById(id int64) (*AssetLatestPrice, error)
	GetAssetLatestPrice(asset *Asset) (*Money, error)
	GetAssetPriceAt(asset *Asset, timestamp time.Time) (*Money, error)
	StoreAssetPrice(asset *Asset, price Money, timestamp time.Time) error
}
