package currencies

type CurrencyService interface {
	GetSourceCurrencies() ([]*Currency, error)
	GetTargetCurrencies() ([]*Currency, error)
}
