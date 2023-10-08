package sqlite

import (
	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	database "github.com/paletas/silvestre.finances/internal/pkg/infrastructure/sqlite/assets"
)

type DatabaseCryptoService struct {
	assetTable *database.CryptoAssetTable
}

func NewDatabaseCryptoService(assetsDb *database.AssetsDB) *DatabaseCryptoService {
	return &DatabaseCryptoService{
		assetTable: assetsDb.CryptoAssetTable,
	}
}

func (s *DatabaseCryptoService) CreateAsset(asset *assets.CryptoAsset) error {
	return s.assetTable.Create(asset)
}

func (s *DatabaseCryptoService) GetAssetByTicker(ticker string) (*assets.CryptoAsset, error) {
	asset, err := s.assetTable.GetByTicker(ticker)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *DatabaseCryptoService) ListAll() ([]*assets.CryptoAsset, error) {
	assets_arr, err := s.assetTable.ListAll()
	if err != nil {
		return nil, err
	}

	cryptoAssets := append([]*assets.CryptoAsset{}, assets_arr...)
	return cryptoAssets, nil
}
