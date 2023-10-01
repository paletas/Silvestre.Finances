package database

import (
	"database/sql"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type StockAssetTable struct {
	db *sql.DB
}

func NewStockAssetTable(db *sql.DB) *StockAssetTable {
	return &StockAssetTable{
		db: db,
	}
}

func (a *StockAssetTable) Create(asset *assets.StockAsset) error {
	_, err := a.db.Exec("EXEC dbo.CreateStockAsset @Name = $1, @Ticker = $2, @Exchange = $3, @Currency = $4", asset.Asset.Name, asset.Ticker, asset.Exchange, asset.Currency)
	if err != nil {
		return err
	}
	return nil
}

func (a *StockAssetTable) GetByTicker(ticker string) (*assets.StockAsset, error) {
	queryResult, err := a.db.Query("EXEC dbo.GetStockAssetByTicker @Ticker = $1", ticker)
	if err != nil {
		return nil, err
	}

	var asset assets.StockAsset
	queryResult.Next()
	err = queryResult.Scan(&asset.Asset.Id, &asset.Asset.Name, &asset.Ticker, &asset.Exchange, &asset.Currency)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (a *StockAssetTable) GetByISIN(isin string) (*assets.StockAsset, error) {
	queryResult, err := a.db.Query("EXEC dbo.GetStockAssetByISIN @ISIN = $1", isin)
	if err != nil {
		return nil, err
	}

	var asset assets.StockAsset
	queryResult.Next()
	err = queryResult.Scan(&asset.Asset.Id, &asset.Asset.Name, &asset.Ticker, &asset.Exchange, &asset.Currency)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (a *StockAssetTable) ListAll() ([]*assets.StockAsset, error) {
	queryResult, err := a.db.Query("EXEC dbo.ListAllStockAssets")
	if err != nil {
		return nil, err
	}

	assets_arr := make([]*assets.StockAsset, 0)
	for queryResult.Next() {
		var asset assets.StockAsset
		err := queryResult.Scan(&asset.Asset.Id, &asset.Asset.Name, &asset.Ticker, &asset.Exchange, &asset.Currency)
		if err != nil {
			return nil, err
		}
		assets_arr = append(assets_arr, &asset)
	}
	return assets_arr, nil
}
