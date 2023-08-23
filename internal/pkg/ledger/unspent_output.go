package ledger

import (
	"errors"
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type UnspentOutput struct {
	Id            string
	TransactionId string
	Date          time.Time
	AssetType     assets.AssetType
	AssetId       string
	Amount        float64
	Spent         bool
	SpentDate     time.Time
}

func createUnspentOutput(transaction_id string, date time.Time, asset_type assets.AssetType, asset_id string, amount float64) UnspentOutput {
	return UnspentOutput{
		TransactionId: transaction_id,
		Date:          date,
		AssetType:     asset_type,
		AssetId:       asset_id,
		Amount:        amount,
		Spent:         false,
		SpentDate:     time.Time{},
	}
}

func (utxo *UnspentOutput) markAsSpent(date time.Time) error {
	if utxo.Spent {
		return errors.New("unspent output already spent")
	}

	utxo.Spent = true
	utxo.SpentDate = date

	return nil
}
