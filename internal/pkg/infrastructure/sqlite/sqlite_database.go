package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paletas/silvestre.finances/internal/pkg/infrastructure/sqlite/tables"
)

type FinancesDb struct {
	db *sql.DB

	StockAssetTable  *tables.StockAssetTable
	CryptoAssetTable *tables.CryptoAssetTable
	AssetPriceTable  *tables.AssetPriceTable
	LedgerTable      *tables.LedgerTable
}

func NewFinancesDb(dbPath string) (*FinancesDb, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	assetsDb := &FinancesDb{
		db: db,

		StockAssetTable:  tables.NewStockAssetTable(db),
		CryptoAssetTable: tables.NewCryptoAssetTable(db),
		AssetPriceTable:  tables.NewAssetPriceTable(db),
		LedgerTable:      tables.NewLedgerTable(db),
	}

	err = assetsDb.applyMigrations()
	if err != nil {
		return nil, err
	}

	return assetsDb, nil
}

func (a *FinancesDb) Disconnect() {
	err := a.db.Close()
	if err != nil {
		panic(err)
	}
}
