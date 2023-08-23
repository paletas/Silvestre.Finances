package ledger

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LedgerDb struct {
	db *mongo.Client

	UnspentOutput UnspentOutputCollection
}

func NewLedgerDB(dbOptions LedgerDbOptions) (*LedgerDb, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbOptions.ConnectionString))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &LedgerDb{
		db:            client,
		UnspentOutput: *NewUnspentOutputCollection(client.Database("ledger").Collection("unspent_outputs")),
	}, nil
}

func (l *LedgerDb) Disconnect() {
	err := l.db.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}
}
