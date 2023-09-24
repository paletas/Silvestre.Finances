package polygon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/currencies"
	"github.com/paletas/silvestre.finances/internal/pkg/feeds"
)

type PolygonService struct {
	PolygonURL    string
	PolygonApiKey string
}

func NewPolygonService(apiKey string) *PolygonService {
	return &PolygonService{
		PolygonURL:    "https://api.polygon.io/",
		PolygonApiKey: apiKey,
	}
}

func (s PolygonService) SearchStock(ticker string) ([]*assets.StockAsset, error) {
	request, err := newPolygonGetRequest(s.PolygonURL, s.PolygonApiKey, fmt.Sprintf("v3/reference/tickers?search=%s&market=stocks", ticker))
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var queryResponse PolygonStockTickersResponse
	json.NewDecoder(response.Body).Decode(&queryResponse)

	var stocks []*assets.StockAsset
	for _, match := range queryResponse.Results {
		stocks = append(stocks, &assets.StockAsset{
			Asset: assets.Asset{
				Name: match.Name,
				Type: assets.StockAssetType,
			},
			Ticker:   match.Ticker,
			Exchange: match.PrimaryExchange,
			Currency: match.Currency,
		})
	}

	return stocks, nil
}

func (s PolygonService) GetStockFromTicker(ticker string) (*assets.StockAsset, error) {
	request, err := newPolygonGetRequest(s.PolygonURL, s.PolygonApiKey, fmt.Sprintf("v3/reference/tickers/%s", ticker))
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var queryResponse PolygonTickerDetailsReponse
	json.NewDecoder(response.Body).Decode(&queryResponse)

	match := queryResponse.Results
	return &assets.StockAsset{
		Asset: assets.Asset{
			Name: match.Name,
			Type: assets.StockAssetType,
		},
		Ticker:   match.Ticker,
		Exchange: match.PrimaryExchange,
		Currency: match.Currency,
	}, nil
}

func (s PolygonService) GetQuoteAtDate(ticker string, date time.Time) (*feeds.StockQuote, error) {
	request, err := newPolygonGetRequest(s.PolygonURL, s.PolygonApiKey, fmt.Sprintf("v1/open-close/%s/%s?adjusted=false", ticker, date.Format("2006-01-02")))
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var getResponse PolygonDailyOpenCloseResponse
	json.NewDecoder(response.Body).Decode(&getResponse)

	return &feeds.StockQuote{
		Ticker: ticker,
		Date:   date,
		Open:   getResponse.Open,
		High:   getResponse.High,
		Low:    getResponse.Low,
		Close:  getResponse.Close,
	}, nil
}

func (s PolygonService) GetSourceCurrencies() ([]*currencies.Currency, error) {
	currenciesChannel := make(chan []*currencies.Currency)
	errorsChannel := make(chan error)
	defer close(currenciesChannel)
	defer close(errorsChannel)

	targetCurrencies, err := s.GetTargetCurrencies()
	if err != nil {
		return nil, err
	}

	for _, targetCurrency := range targetCurrencies {
		go func(targetCurrency *currencies.Currency) {
			currencies, err := s.GetSourceCurrency(targetCurrency)
			if err != nil {
				errorsChannel <- err
				return
			}

			currenciesChannel <- currencies
		}(targetCurrency)
	}

	var found []*currencies.Currency
	for i := 0; i < len(targetCurrencies); i++ {
		select {
		case currencies := <-currenciesChannel:
			found = append(found, currencies...)
		case err := <-errorsChannel:
			return nil, err
		}
	}
	return found, nil
}

func (s PolygonService) GetSourceCurrency(targetCurrency *currencies.Currency) ([]*currencies.Currency, error) {
	request, err := newPolygonGetRequest(s.PolygonURL, s.PolygonApiKey, fmt.Sprintf("v3/reference/tickers?search=C:%s&market=fx", targetCurrency.Code))
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var queryResponse PolygonCurrencyTickersResponse
	json.NewDecoder(response.Body).Decode(&queryResponse)

	var found []*currencies.Currency
	for _, match := range queryResponse.Results {
		found = append(found, &currencies.Currency{
			Name: match.CurrencyName,
			Code: match.CurrencySymbol,
		})
	}

	return found, nil
}

func (s PolygonService) GetTargetCurrencies() ([]*currencies.Currency, error) {
	return []*currencies.Currency{
		{
			Name: "United States Dollar",
			Code: "USD",
		},
		{
			Name: "Euro",
			Code: "EUR",
		},
	}, nil
}

func newPolygonGetRequest(basePath string, apiKey string, path string) (*http.Request, error) {
	absolute_path := basePath + path
	req, err := http.NewRequest("GET", absolute_path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)

	return req, nil
}

type PolygonStockTickersResponse struct {
	Count     int    `json:"count"`
	NextUrl   string `json:"next_url"`
	RequestId string `json:"request_id"`
	Results   []struct {
		Active          bool   `json:"active"`
		CIK             string `json:"cik"`
		CompositeFIGI   string `json:"composite_figi"`
		Currency        string `json:"currency_name"`
		LastUpdatedUTC  string `json:"last_updated_utc"`
		Locale          string `json:"locale"`
		Market          string `json:"market"`
		Name            string `json:"name"`
		PrimaryExchange string `json:"primary_exchange"`
		ShareClassFIGI  string `json:"share_class_figi"`
		Ticker          string `json:"ticker"`
		Type            string `json:"type"`
	} `json:"results"`
	Status string `json:"status"`
}

type PolygonDailyOpenCloseResponse struct {
	AfterHours float64 `json:"afterHours"`
	Close      float64 `json:"close"`
	From       string  `json:"from"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Open       float64 `json:"open"`
	PreMarket  float64 `json:"preMarket"`
	Status     string  `json:"status"`
	Symbol     string  `json:"symbol"`
	Volume     int     `json:"volume"`
}

type PolygonTickerDetailsReponse struct {
	RequestId string `json:"request_id"`
	Results   struct {
		Active  bool `json:"active"`
		Address struct {
			Address1   string `json:"address1"`
			City       string `json:"city"`
			PostalCode string `json:"postalCode"`
			State      string `json:"state"`
		}
		Branding struct {
			IconUrl string `json:"icon"`
			LogoUrl string `json:"logo"`
		}
		CIK                         string `json:"cik"`
		CompositeFIGI               string `json:"composite_figi"`
		Currency                    string `json:"currency_name"`
		Description                 string `json:"description"`
		Homepage                    string `json:"homepage_url"`
		ListDate                    string `json:"list_date"`
		Locale                      string `json:"locale"`
		Market                      string `json:"market"`
		MarketCap                   int64  `json:"market_cap"`
		Name                        string `json:"name"`
		PhoneNumber                 string `json:"phone"`
		PrimaryExchange             string `json:"primary_exchange"`
		RoundLot                    int    `json:"round_lot"`
		ShareClassFIGI              string `json:"share_class_figi"`
		ShareClassSharesOutstanding int64  `json:"share_class_shares_outstanding"`
		SICCode                     string `json:"sic_code"`
		SICDDescription             string `json:"sic_description"`
		Ticker                      string `json:"ticker"`
		TickerRoot                  string `json:"ticker_root"`
		TotalEmployees              int64  `json:"total_employees"`
		Type                        string `json:"type"`
		WeightedSharesOutstanding   int64  `json:"weighted_shares_outstanding"`
	} `json:"results"`
	Status string `json:"status"`
}

type PolygonCurrencyTickersResponse struct {
	Count     int    `json:"count"`
	NextUrl   string `json:"next_url"`
	RequestId string `json:"request_id"`
	Results   []struct {
		Ticker             string `json:"ticker"`
		Name               string `json:"name"`
		Market             string `json:"market"`
		Locale             string `json:"locale"`
		Active             bool   `json:"active"`
		CurrencySymbol     string `json:"currency_symbol"`
		CurrencyName       string `json:"currency_name"`
		BaseCurrencySymbol string `json:"base_currency_symbol"`
		BaseCurrencyName   string `json:"base_currency_name"`
		LastUpdatedUTC     string `json:"last_updated_utc"`
	} `json:"results"`
	Status string `json:"status"`
}
