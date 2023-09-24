package database

import (
	"database/sql"
)

type AssetsDB struct {
	db *sql.DB

	StockAssetTable  *StockAssetTable
	CryptoAssetTable *CryptoAssetTable
}

func NewAssetsDB(db *sql.DB) *AssetsDB {
	return &AssetsDB{
		db: db,

		StockAssetTable:  NewStockAssetTable(db),
		CryptoAssetTable: NewCryptoAssetTable(db),
	}
}

func (a *AssetsDB) Disconnect() {
	err := a.db.Close()
	if err != nil {
		panic(err)
	}
}
