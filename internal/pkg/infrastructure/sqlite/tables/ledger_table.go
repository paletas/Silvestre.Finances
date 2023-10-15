package tables

import (
	"context"
	"database/sql"
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type LedgerTable struct {
	db *sql.DB
}

func NewLedgerTable(db *sql.DB) *LedgerTable {
	return &LedgerTable{
		db: db,
	}
}

func (l *LedgerTable) AddUnspentOutput(
	transaction_id string,
	exchange string,
	date time.Time,
	asset_type string,
	asset_id int64,
	amount float64,
	costBasis float64,
	costBasisCurrency string,
	fees float64,
	feesCurrency string) error {

	conn, err := l.db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO ledger(transaction_id, exchange, date, asset_type, asset_id, amount, cost_basis, cost_basis_currency, fees, fees_currency) VALUES(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(transaction_id, exchange, date, asset_type, asset_id, amount, costBasis, costBasisCurrency, fees, feesCurrency)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (l *LedgerTable) SpendOutput(transaction_id string, date time.Time, fees float64) error {
	conn, err := l.db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE ledger SET spent = 1, spent_at = ?, fees = fees + ? WHERE transaction_id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(date, fees, transaction_id)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (l *LedgerTable) GetUnspentOutputs() ([]ledger.UnspentOutput, error) {
	conn, err := l.db.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(context.Background(), "SELECT transaction_id, date, asset_id, amount, cost_basis, cost_basis_currency, fees, fees_currency FROM ledger WHERE spent = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	utxos := make([]ledger.UnspentOutput, 0)
	for rows.Next() {
		var utxo ledger.UnspentOutput
		err := rows.Scan(&utxo.TransactionId, &utxo.Date, &utxo.AssetId, &utxo.Amount, &utxo.CostBasis.Amount, &utxo.CostBasis.Currency, &utxo.Fees.Amount, &utxo.Fees.Currency)
		if err != nil {
			return nil, err
		}

		utxos = append(utxos, utxo)
	}

	return utxos, nil
}

func (l *LedgerTable) GetUnspentOutputsByAssetType(assetType string) ([]ledger.UnspentOutput, error) {
	conn, err := l.db.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(context.Background(), "SELECT transaction_id, date, asset_id, amount, cost_basis, cost_basis_currency, fees, fees_currency FROM ledger WHERE spent = 0 AND asset_type=?", assetType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	utxos := make([]ledger.UnspentOutput, 0)
	for rows.Next() {
		var utxo ledger.UnspentOutput
		err := rows.Scan(&utxo.TransactionId, &utxo.Date, &utxo.AssetId, &utxo.Amount, &utxo.CostBasis.Amount, &utxo.CostBasis.Currency, &utxo.Fees.Amount, &utxo.Fees.Currency)
		if err != nil {
			return nil, err
		}

		utxos = append(utxos, utxo)
	}

	return utxos, nil
}
