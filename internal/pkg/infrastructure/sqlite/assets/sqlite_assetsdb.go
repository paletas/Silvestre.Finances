package assets

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type AssetsDB struct {
	db *sql.DB

	StockAssetTable  *StockAssetTable
	CryptoAssetTable *CryptoAssetTable
	AssetPriceTable  *AssetPriceTable
}

func NewAssetsDb(dbPath string) (*AssetsDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	assetsDb := &AssetsDB{
		db: db,

		StockAssetTable:  NewStockAssetTable(db),
		CryptoAssetTable: NewCryptoAssetTable(db),
		AssetPriceTable:  NewAssetPriceTable(db),
	}

	err = assetsDb.applyMigrations()
	if err != nil {
		return nil, err
	}

	return assetsDb, nil
}

func (a *AssetsDB) Disconnect() {
	err := a.db.Close()
	if err != nil {
		panic(err)
	}
}
