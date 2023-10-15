package ledger

import (
	"errors"
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type UnspentOutput struct {
	Id            string
	TransactionId string
	Exchange      string
	Date          time.Time
	AssetId       int64
	Amount        float64
	CostBasis     assets.Money
	Fees          assets.Money
	Spent         bool
	SpentDate     time.Time
	SpendFees     float64
}

func CreateUnspentOutput(
	transaction_id string,
	exchange string,
	date time.Time,
	asset_type assets.AssetType,
	asset_id int64,
	amount float64,
	costBasis assets.Money,
	fees assets.Money) UnspentOutput {

	return UnspentOutput{
		TransactionId: transaction_id,
		Exchange:      exchange,
		Date:          date,
		AssetId:       asset_id,
		Amount:        amount,
		CostBasis:     costBasis,
		Fees:          fees,
		Spent:         false,
		SpentDate:     time.Time{},
	}
}

func (utxo *UnspentOutput) MarkAsSpent(date time.Time, fees float64) error {
	if utxo.Spent {
		return errors.New("unspent output already spent")
	}

	utxo.Spent = true
	utxo.SpentDate = date
	utxo.SpendFees = fees

	return nil
}
