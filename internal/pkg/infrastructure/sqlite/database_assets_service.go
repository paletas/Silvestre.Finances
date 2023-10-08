package sqlite

import (
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	database "github.com/paletas/silvestre.finances/internal/pkg/infrastructure/sqlite/assets"
)

type DatabaseAssetsService struct {
	assetPriceTable *database.AssetPriceTable
}

func NewDatabaseAssetsService(assetsDb *database.AssetsDB) *DatabaseAssetsService {
	return &DatabaseAssetsService{
		assetPriceTable: assetsDb.AssetPriceTable,
	}
}

func (s *DatabaseAssetsService) GetAssetById(id int64) (*assets.Asset, error) {
	asset, err := s.assetPriceTable.GetAssetById(id)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *DatabaseAssetsService) GetAssetLatestPrice(asset *assets.Asset) (*assets.Money, error) {
	price, currency, err := s.assetPriceTable.GetLatestPrice(asset)
	if err != nil {
		return nil, err
	}
	return &assets.Money{Amount: price, Currency: currency}, nil
}

func (s *DatabaseAssetsService) GetAssetPriceAt(asset *assets.Asset, timestamp time.Time) (*assets.Money, error) {
	price, currency, err := s.assetPriceTable.GetPriceAt(asset, timestamp)
	if err != nil {
		return nil, err
	}
	return &assets.Money{Amount: price, Currency: currency}, nil
}

func (s *DatabaseAssetsService) StoreAssetPrice(asset *assets.Asset, price assets.Money, timestamp time.Time) error {
	return s.assetPriceTable.RegisterPrice(asset, price.Amount, price.Currency, timestamp)
}
