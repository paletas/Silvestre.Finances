package tables

import (
	"context"
	"database/sql"
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type AssetPriceTable struct {
	db *sql.DB
}

func NewAssetPriceTable(db *sql.DB) *AssetPriceTable {
	return &AssetPriceTable{db: db}
}

func (t *AssetPriceTable) RegisterPrice(asset *assets.Asset, price float64, currency string, date time.Time) error {
	conn, err := t.db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(context.Background(), `
		INSERT INTO AssetPrice (AssetID, Date, Price, Currency)
		VALUES (?, ?, ?)`, asset.Id, date, price, currency)
	if err != nil {
		return err
	}

	return nil
}

func (t *AssetPriceTable) GetPriceAt(asset *assets.Asset, date time.Time) (float64, string, error) {
	conn, err := t.db.Conn(context.Background())
	if err != nil {
		return 0.0, "", err
	}
	defer conn.Close()

	var price float64
	var currency string
	err = conn.QueryRowContext(context.Background(), `
		SELECT Price, Currency
		FROM AssetPrice
		WHERE AssetID = ?
		AND Date = ?`, asset.Id, date).Scan(&price, &currency)
	if err != nil {
		return 0.0, "", err
	}

	return price, currency, nil
}

func (t *AssetPriceTable) GetLatestPrice(asset *assets.Asset) (float64, string, error) {
	conn, err := t.db.Conn(context.Background())
	if err != nil {
		return 0.0, "", err
	}
	defer conn.Close()

	var price float64
	var currency string
	err = conn.QueryRowContext(context.Background(), `
		SELECT Price, Currency
		FROM AssetPrice
		WHERE AssetID = ?
		ORDER BY Date DESC
		LIMIT 1`, asset.Id).Scan(&price, &currency)
	if err != nil {
		return 0.0, "", err
	}

	return price, currency, nil
}

func (t *AssetPriceTable) GetAssetById(id int64) (*assets.Asset, error) {
	conn, err := t.db.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var asset assets.Asset
	err = conn.QueryRowContext(context.Background(), `
		SELECT A.ID, A.Name, A.AssetType
		FROM Asset A
		WHERE A.ID = ?`, id).Scan(&asset.Id, &asset.Name, &asset.Type)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}
