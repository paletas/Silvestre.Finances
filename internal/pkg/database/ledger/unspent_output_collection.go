package ledger

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/paletas/silvestre.finances/internal/pkg/database/ledger/models"
	"github.com/paletas/silvestre.finances/internal/pkg/ledger"
)

type UnspentOutputCollection struct {
	collection *mongo.Collection
}

func NewUnspentOutputCollection(collection *mongo.Collection) *UnspentOutputCollection {
	return &UnspentOutputCollection{collection: collection}
}

func (uc *UnspentOutputCollection) InsertOne(utxo *ledger.UnspentOutput) (string, error) {
	db_utxo, err := models.MapToDatabaseModel(utxo)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	result, err := uc.collection.InsertOne(context.Background(), db_utxo)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	id, convertion_worked := result.InsertedID.(primitive.ObjectID)
	if convertion_worked == false {
		return "", errors.New("Could not convert inserted ID to string")
	}

	return id.Hex(), nil
}
