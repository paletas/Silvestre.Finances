package sqlite

import (
	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	database "github.com/paletas/silvestre.finances/internal/pkg/infrastructure/sqlite/assets"
)

type DatabaseStocksService struct {
	assetTable *database.StockAssetTable
}

func NewDatabaseStocksService(assetsDb *database.AssetsDB) *DatabaseStocksService {
	return &DatabaseStocksService{
		assetTable: assetsDb.StockAssetTable,
	}
}

func (s *DatabaseStocksService) CreateAsset(asset *assets.StockAsset) error {
	return s.assetTable.Create(asset)
}

func (s *DatabaseStocksService) GetAssetByTicker(ticker string) (*assets.StockAsset, error) {
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
