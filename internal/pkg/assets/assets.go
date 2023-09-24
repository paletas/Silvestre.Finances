package assets

type AssetType string

const (
	StockAssetType  AssetType = "STOCK"
	CryptoAssetType AssetType = "CRYPTO"
)

type Asset struct {
	Id   string
	Name string
	Type AssetType
}
