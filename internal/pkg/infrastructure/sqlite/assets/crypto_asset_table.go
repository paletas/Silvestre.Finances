package assets

import (
	"context"
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
	conn, err := a.db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(context.Background(), `
		PRAGMA temp_store = 2;
		CREATE TEMP TABLE Variables(AssetID INTEGER);

		INSERT INTO Asset (AssetType, Name)
		VALUES ('Crypto', ?);

		INSERT INTO Variables(AssetID)
		VALUES (last_insert_rowid());

		INSERT INTO CryptoAsset (ID, Ticker)
		SELECT AssetID, ?
		FROM Variables;

		SELECT AssetID FROM Variables;`, asset.Asset.Name, asset.Ticker)
	if err != nil {
		return err
	}
	return nil
}

func (a *CryptoAssetTable) GetByTicker(ticker string) (*assets.CryptoAsset, error) {
	conn, err := a.db.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	queryResult, err := conn.QueryContext(context.Background(), `
		SELECT A.ID, A.Name, C.Ticker
		FROM Asset A
		INNER JOIN CryptoAsset C ON C.ID = A.ID
		WHERE C.Ticker = ?`, ticker)
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
	conn, err := a.db.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	queryResult, err := conn.QueryContext(context.Background(), `
		SELECT A.ID, A.Name, C.Ticker
		FROM Asset A
		INNER JOIN CryptoAsset C ON C.ID = A.ID`)
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
