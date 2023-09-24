package assets

type CryptoAssetsService interface {
	CreateAsset(asset *CryptoAsset) error
	GetAssetByTicker(ticker string) (*CryptoAsset, error)
	ListAll() ([]*CryptoAsset, error)
}
