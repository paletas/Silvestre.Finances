package assets

import (
	"context"
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
	conn, err := a.db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(context.Background(), `
		PRAGMA temp_store = 2;
		CREATE TEMP TABLE Variables(AssetID INTEGER);

		INSERT INTO Asset (AssetType, Name)
		VALUES ('Stock', ?);

		INSERT INTO Variables(AssetID)
		VALUES (last_insert_rowid());

		INSERT INTO StockAsset (ID, Ticker, Exchange, Currency)
		SELECT AssetID, ?, ?, ?
		FROM Variables;

		SELECT AssetID FROM Variables;`, asset.Asset.Name, asset.Ticker, asset.Exchange, asset.Currency)
	if err != nil {
		return err
	}

	return nil
}

func (a *StockAssetTable) GetByTicker(ticker string) (*assets.StockAsset, error) {
	conn, err := a.db.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	queryResult, err := conn.QueryContext(context.Background(), `
		SELECT A.ID, A.Name, S.Ticker, S.Exchange, S.Currency
		FROM Asset A
		INNER JOIN StockAsset S ON S.ID = A.ID
		WHERE S.Ticker = ?`, ticker)
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
	conn, err := a.db.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	queryResult, err := conn.QueryContext(context.Background(), `
		SELECT A.ID, A.Name, S.Ticker, S.Exchange, S.Currency
		FROM Asset A
		INNER JOIN StockAsset S ON S.ID = A.ID`)
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
