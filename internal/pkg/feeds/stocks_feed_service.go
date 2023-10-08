package feeds

import (
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
)

type StocksFeedService interface {
	SearchStock(ticker string) ([]*assets.StockAsset, error)
	GetStockFromTicker(ticker string) (*assets.StockAsset, error)

	GetQuoteAtDate(ticker string, date time.Time) (*StockQuote, error)
}
