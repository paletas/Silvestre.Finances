package inmemory

import (
	"errors"
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type MemoryLedger struct {
	unspent map[string]ledger.UnspentOutput
	spent   []ledger.UnspentOutput
}

func CreateMemoryLedger() MemoryLedger {
	return MemoryLedger{
		unspent: make(map[string]ledger.UnspentOutput),
		spent:   make([]ledger.UnspentOutput, 0),
	}
}

func (l MemoryLedger) AddUnspentOutput(
	transaction_id string,
	date time.Time,
	asset_type assets.AssetType,
	asset_id string,
	amount float64,
	costBasis float64,
	fees float64) error {

	utxo := ledger.CreateUnspentOutput(transaction_id, date, asset_type, asset_id, amount, costBasis, fees)
	l.unspent[transaction_id] = utxo

	return nil
}

func (l MemoryLedger) SpendOutput(transaction_id string, date time.Time, fees float64) error {
	utxo, ok := l.unspent[transaction_id]
	if !ok {
		return errors.New("unspent output not found")
	}

	err := utxo.MarkAsSpent(date, fees)
	if err != nil {
		return err
	}

	l.spent = append(l.spent, utxo)
	delete(l.unspent, transaction_id)

	return nil
}

func (l MemoryLedger) GetTotalWealth() (float64, error) {
	total := 0.0
	for _, utxo := range l.unspent {
		total += utxo.Amount
	}

	return total, nil
}

func (l MemoryLedger) GetUnspentOutputs() ([]ledger.UnspentOutput, error) {
	utxos := make([]ledger.UnspentOutput, 0)
	for _, utxo := range l.unspent {
		utxos = append(utxos, utxo)
	}

	return utxos, nil
}

func (l MemoryLedger) GetUnspentOutputsByAssetType(assetType assets.AssetType) ([]ledger.UnspentOutput, error) {
	utxos := make([]ledger.UnspentOutput, 0)
	for _, utxo := range l.unspent {
		if utxo.AssetType == assetType {
			utxos = append(utxos, utxo)
		}
	}

	return utxos, nil
}
