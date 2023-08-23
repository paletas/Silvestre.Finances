package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/paletas/silvestre.finances/internal/pkg/assets"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type unspentOutput struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	TransactionId string             `bson:"transaction_id"`
	Date          time.Time          `bson:"date"`
	AssetType     string             `bson:"asset_type"`
	AssetId       string             `bson:"asset_id"`
	Amount        float64            `bson:"amount"`
	Spent         bool               `bson:"spent"`
	SpentDate     time.Time          `bson:"spent_date"`
}

func MapToDatabaseModel(utxo *ledger.UnspentOutput) (*unspentOutput, error) {
	var objid primitive.ObjectID
	if utxo.Id == "" {
		objid = primitive.NewObjectID()
	} else {
		var err error
		objid, err = primitive.ObjectIDFromHex(utxo.Id)
		if err != nil {
			return nil, err
		}
	}

	return &unspentOutput{
		Id:            objid,
		TransactionId: utxo.TransactionId,
		Date:          utxo.Date,
		AssetType:     string(utxo.AssetType),
		AssetId:       utxo.AssetId,
		Amount:        utxo.Amount,
		Spent:         utxo.Spent,
		SpentDate:     utxo.SpentDate,
	}, nil
}

func (utxo *unspentOutput) MapToModel() *ledger.UnspentOutput {
	return &ledger.UnspentOutput{
		Id:            utxo.Id.Hex(),
		TransactionId: utxo.TransactionId,
		Date:          utxo.Date,
		AssetType:     assets.AssetType(utxo.AssetType),
		AssetId:       utxo.AssetId,
		Amount:        utxo.Amount,
		Spent:         utxo.Spent,
		SpentDate:     utxo.SpentDate,
	}
}
