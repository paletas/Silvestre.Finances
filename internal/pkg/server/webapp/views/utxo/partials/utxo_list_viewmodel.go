package partials

type UnspentOutputListViewModel struct {
	UnspentOutputs []UnspentOutputViewModel
}

type UnspentOutputViewModel struct {
	Id            string
	TransactionId string
	Exchange      string
	Date          string
	AssetId       int64
	AssetName     string
	Value         float64
}
