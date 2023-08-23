package ledger

import (
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type Ledger interface {
	AddUnspentOutput(transaction_id string, date time.Time, asset_type assets.AssetType, asset_id string, amount float64)
	SpendOutput(transaction_id string, date time.Time) error
}
