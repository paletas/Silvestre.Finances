package sqlite

import (
	"log"
	"time"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	database "github.com/paletas/silvestre.finances/internal/pkg/infrastructure/sqlite/ledger"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type DatabaseLedgerService struct {
	ledgerTable   *database.LedgerTable
	assetsService assets.AssetsService
}

func NewDatabaseLedgerService(ledgerDb *database.LedgerDb, assetsService assets.AssetsService) *DatabaseLedgerService {
	return &DatabaseLedgerService{
		ledgerTable:   ledgerDb.LedgerTable,
		assetsService: assetsService,
	}
}

func (s DatabaseLedgerService) AddUnspentOutput(
	transaction_id string,
	exchange string,
	date time.Time,
	asset_type assets.AssetType,
	asset_id int64,
	amount float64,
	costBasis assets.Money,
	fees assets.Money) error {

	return s.ledgerTable.AddUnspentOutput(
		transaction_id,
		exchange, date,
		string(asset_type),
		asset_id,
		amount,
		costBasis.Amount,
		costBasis.Currency,
		fees.Amount,
		fees.Currency)
}

func (s DatabaseLedgerService) SpendOutput(transaction_id string, date time.Time, fees float64) error {
	return s.ledgerTable.SpendOutput(transaction_id, date, fees)
}

func (s DatabaseLedgerService) GetTotalWealth() (float64, error) {
	utxos, err := s.ledgerTable.GetUnspentOutputs()
	if err != nil {
		return 0.0, err
	}

	total := 0.0
	for _, utxo := range utxos {
		asset, err := s.assetsService.GetAssetById(utxo.AssetId)
		if err != nil {
			log.Printf("Error getting asset by id: %v", err)
			continue
		}

		assetPrice, err := s.assetsService.GetAssetLatestPrice(asset)
		if err != nil {
			log.Printf("Error getting asset latest price: %v", err)
			continue
		}

		total += utxo.Amount * assetPrice.Amount
	}
	return total, nil
}

func (s DatabaseLedgerService) GetUnspentOutputs() ([]ledger.UnspentOutput, error) {
	return s.ledgerTable.GetUnspentOutputs()
}

func (s DatabaseLedgerService) GetUnspentOutputsByAssetType(assetType assets.AssetType) ([]ledger.UnspentOutput, error) {
	return s.ledgerTable.GetUnspentOutputsByAssetType(string(assetType))
}
