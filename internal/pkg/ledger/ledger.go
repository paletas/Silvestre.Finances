package ledger

import (
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type Ledger interface {
	AddUnspentOutput(transaction_id string, date time.Time, asset_type assets.AssetType, asset_id string, amount float64, costBasis float64, fees float64) error
	SpendOutput(transaction_id string, date time.Time, fees float64) error
	GetTotalWealth() (float64, error)
	GetUnspentOutputs() ([]UnspentOutput, error)
	GetUnspentOutputsByAssetType(assetType assets.AssetType) ([]UnspentOutput, error)
}
