package ledger

import (
	"errors"
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type MemoryLedger struct {
	unspent map[string]UnspentOutput
	spent   []UnspentOutput
}

func CreateMemoryLedger() MemoryLedger {
	return MemoryLedger{
		unspent: make(map[string]UnspentOutput),
		spent:   make([]UnspentOutput, 0),
	}
}

func (l *MemoryLedger) AddUnspentOutput(transaction_id string, date time.Time, asset_type assets.AssetType, asset_id string, amount float64) {
	utxo := createUnspentOutput(transaction_id, date, asset_type, asset_id, amount)
	l.unspent[transaction_id] = utxo
}

func (l *MemoryLedger) SpendOutput(transaction_id string, date time.Time) error {
	utxo, ok := l.unspent[transaction_id]
	if !ok {
		return errors.New("unspent output not found")
	}

	err := utxo.markAsSpent(date)
	if err != nil {
		return err
	}

	l.spent = append(l.spent, utxo)

	return nil
}
