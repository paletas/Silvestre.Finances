package database

import (
	"database/sql"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type CryptoAssetTable struct {
	db *sql.DB
}

func NewCryptoAssetTable(db *sql.DB) *CryptoAssetTable {
	return &CryptoAssetTable{
		db: db,
	}
}

func (a *CryptoAssetTable) Create(asset *assets.CryptoAsset) error {
	_, err := a.db.Exec("EXEC dbo.CreateCryptoAsset @Name = $1, @Ticker = $2", asset.Asset.Name, asset.Ticker)
	if err != nil {
		return err
	}
	return nil
}

func (a *CryptoAssetTable) GetByTicker(ticker string) (*assets.CryptoAsset, error) {
	queryResult, err := a.db.Query("EXEC dbo.GetCryptoAssetByTicker @Ticker = $1", ticker)
	if err != nil {
		return nil, err
	}

	var asset assets.CryptoAsset
	queryResult.Next()
	err = queryResult.Scan(&asset.Asset.Id, &asset.Asset.Name, &asset.Ticker)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (a *CryptoAssetTable) ListAll() ([]*assets.CryptoAsset, error) {
	queryResult, err := a.db.Query("EXEC dbo.ListAllCryptoAssets")
	if err != nil {
		return nil, err
	}

	assets_arr := make([]*assets.CryptoAsset, 0)
	for queryResult.Next() {
		var asset assets.CryptoAsset
		err := queryResult.Scan(&asset.Asset.Id, &asset.Asset.Name, &asset.Ticker)
		if err != nil {
			return nil, err
		}
		assets_arr = append(assets_arr, &asset)
	}
	return assets_arr, nil
}
