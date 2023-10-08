package assets

import (
	"time"
)

type AssetPrice struct {
	AssetID int64
	Price   Money
	Date    time.Time
}
