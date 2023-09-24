package feeds

import (
	"time"
)

type StockQuote struct {
	Ticker string
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
}
