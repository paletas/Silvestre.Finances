package assets

type StockAssetsService interface {
	CreateAsset(asset *StockAsset) error
	GetAssetByTicker(ticker string) (*StockAsset, error)
	ListAll() ([]*StockAsset, error)
}
