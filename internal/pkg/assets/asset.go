package assets

type AssetType string

const (
	StockAssetType  AssetType = "STOCK"
	CryptoAssetType AssetType = "CRYPTO"
)

type Asset struct {
	Id   int64
	Name string
	Type AssetType
}
