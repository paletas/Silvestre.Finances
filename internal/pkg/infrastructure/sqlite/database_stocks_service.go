package sqlite

import (
	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type DatabaseStocksService struct {
	db *FinancesDb
}

func NewDatabaseStocksService(db *FinancesDb) *DatabaseStocksService {
	return &DatabaseStocksService{
		db: db,
	}
}

func (s *DatabaseStocksService) CreateAsset(asset *assets.StockAsset) error {
	return s.db.StockAssetTable.Create(asset)
}

func (s *DatabaseStocksService) GetAssetByTicker(ticker string) (*assets.StockAsset, error) {
	asset, err := s.db.StockAssetTable.GetByTicker(ticker)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *DatabaseStocksService) ListAll() ([]*assets.StockAsset, error) {
	assets_arr, err := s.db.StockAssetTable.ListAll()
	if err != nil {
		return nil, err
	}

	stockAssets := append([]*assets.StockAsset{}, assets_arr...)
	return stockAssets, nil
}
