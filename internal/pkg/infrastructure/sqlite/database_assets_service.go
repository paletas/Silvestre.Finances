package sqlite

import (
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type DatabaseAssetsService struct {
	db *FinancesDb
}

func NewDatabaseAssetsService(db *FinancesDb) *DatabaseAssetsService {
	return &DatabaseAssetsService{
		db: db,
	}
}

func (s *DatabaseAssetsService) GetAssetById(id int64) (*assets.Asset, error) {
	asset, err := s.db.AssetPriceTable.GetAssetById(id)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *DatabaseAssetsService) GetAssetLatestPriceById(id int64) (*assets.AssetLatestPrice, error) {
	asset, err := s.db.AssetPriceTable.GetAssetById(id)
	if err != nil {
		return nil, err
	}
	price, currency, err := s.db.AssetPriceTable.GetLatestPrice(asset)
	if err != nil {
		return nil, err
	}
	return &assets.AssetLatestPrice{
		Asset: *asset,
		LatestPrice: assets.Money{
			Amount:   price,
			Currency: currency}}, nil
}

func (s *DatabaseAssetsService) GetAssetLatestPrice(asset *assets.Asset) (*assets.Money, error) {
	price, currency, err := s.db.AssetPriceTable.GetLatestPrice(asset)
	if err != nil {
		return nil, err
	}
	return &assets.Money{Amount: price, Currency: currency}, nil
}

func (s *DatabaseAssetsService) GetAssetPriceAt(asset *assets.Asset, timestamp time.Time) (*assets.Money, error) {
	price, currency, err := s.db.AssetPriceTable.GetPriceAt(asset, timestamp)
	if err != nil {
		return nil, err
	}
	return &assets.Money{Amount: price, Currency: currency}, nil
}

func (s *DatabaseAssetsService) StoreAssetPrice(asset *assets.Asset, price assets.Money, timestamp time.Time) error {
	return s.db.AssetPriceTable.RegisterPrice(asset, price.Amount, price.Currency, timestamp)
}
