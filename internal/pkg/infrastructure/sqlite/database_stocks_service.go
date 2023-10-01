package database

import (
	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	database "github.com/paletas/silvestre.finances/internal/pkg/infrastructure/sqlite/assets"
)

type DatabaseStocksService struct {
	assetTable *database.StockAssetTable
}

func NewDatabaseStocksService(assetTable *database.StockAssetTable) *DatabaseStocksService {
	return &DatabaseStocksService{
		assetTable: assetTable,
	}
}

func (s *DatabaseStocksService) Create(asset *assets.StockAsset) error {
	return s.assetTable.Create(asset)
}

func (s *DatabaseStocksService) GetByTicker(ticker string) (*assets.StockAsset, error) {
	asset, err := s.assetTable.GetByTicker(ticker)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *DatabaseStocksService) ListAll() ([]*assets.StockAsset, error) {
	assets_arr, err := s.assetTable.ListAll()
	if err != nil {
		return nil, err
	}

	stockAssets := append([]*assets.StockAsset{}, assets_arr...)
	return stockAssets, nil
}
