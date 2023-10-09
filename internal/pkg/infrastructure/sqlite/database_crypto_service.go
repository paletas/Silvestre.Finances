package sqlite

import (
	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type DatabaseCryptoService struct {
	db *FinancesDb
}

func NewDatabaseCryptoService(db *FinancesDb) *DatabaseCryptoService {
	return &DatabaseCryptoService{
		db: db,
	}
}

func (s *DatabaseCryptoService) CreateAsset(asset *assets.CryptoAsset) error {
	return s.db.CryptoAssetTable.Create(asset)
}

func (s *DatabaseCryptoService) GetAssetByTicker(ticker string) (*assets.CryptoAsset, error) {
	asset, err := s.db.CryptoAssetTable.GetByTicker(ticker)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *DatabaseCryptoService) ListAll() ([]*assets.CryptoAsset, error) {
	assets_arr, err := s.db.CryptoAssetTable.ListAll()
	if err != nil {
		return nil, err
	}

	cryptoAssets := append([]*assets.CryptoAsset{}, assets_arr...)
	return cryptoAssets, nil
}
