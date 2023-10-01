package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type SqliteLedger struct {
	db *sql.DB
}

func NewSqliteLedger(dbPath string) (*SqliteLedger, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &SqliteLedger{
		db: db,
	}, nil
}

func (l *SqliteLedger) Disconnect() {
	err := l.db.Close()
	if err != nil {
		panic(err)
	}
}

func (l *SqliteLedger) AddUnspentOutput(
	transaction_id string,
	date string,
	asset_type string,
	asset_id string,
	amount float64,
	costBasis float64,
	fees float64) error {

	stmt, err := l.db.Prepare("INSERT INTO unspent_outputs(transaction_id, date, asset_type, asset_id, amount, cost_basis, fees) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(transaction_id, date, asset_type, asset_id, amount, costBasis, fees)
	if err != nil {
		return err
	}

	return nil
}

func (l *SqliteLedger) SpendOutput(transaction_id string, date string, fees float64) error {
	stmt, err := l.db.Prepare("INSERT INTO spent_outputs(transaction_id, date, fees) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(transaction_id, date, fees)
	if err != nil {
		return err
	}

	return nil
}

func (l *SqliteLedger) GetTotalWealth() (float64, error) {
	var total float64
	err := l.db.QueryRow("SELECT SUM(amount) FROM unspent_outputs").Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (l *SqliteLedger) GetUnspentOutputs() ([]ledger.UnspentOutput, error) {
	rows, err := l.db.Query("SELECT transaction_id, date, asset_type, asset_id, amount, cost_basis, fees FROM unspent_outputs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	utxos := make([]ledger.UnspentOutput, 0)
	for rows.Next() {
		var utxo ledger.UnspentOutput
		err := rows.Scan(&utxo.TransactionId, &utxo.Date, &utxo.AssetType, &utxo.AssetId, &utxo.Amount, &utxo.CostBasis, &utxo.Fees)
		if err != nil {
			return nil, err
		}

		utxos = append(utxos, utxo)
	}

	return utxos, nil
}

func (l *SqliteLedger) GetUnspentOutputsByAssetType(assetType string) ([]ledger.UnspentOutput, error) {
	rows, err := l.db.Query("SELECT transaction_id, date, asset_type, asset_id, amount, cost_basis, fees FROM unspent_outputs WHERE asset_type=?", assetType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	utxos := make([]ledger.UnspentOutput, 0)
	for rows.Next() {
		var utxo ledger.UnspentOutput
		err := rows.Scan(&utxo.TransactionId, &utxo.Date, &utxo.AssetType, &utxo.AssetId, &utxo.Amount, &utxo.CostBasis, &utxo.Fees)
		if err != nil {
			return nil, err
		}

		utxos = append(utxos, utxo)
	}

	return utxos, nil
}
