package assets

type StockAsset struct {
	Asset
	Ticker   string
	Exchange string
	Currency string
}
