package ledger

type LedgerDbOptions struct {
	ConnectionString string
	DatabaseName     string
	CollectionName   string
}

func NewLedgerDbOptions(connectionString string) *LedgerDbOptions {
	return NewLedgerDbOptionWithDatabaseName(connectionString, "ledgerdb")
}

func NewLedgerDbOptionWithDatabaseName(connectionString string, databaseName string) *LedgerDbOptions {
	return &LedgerDbOptions{
		ConnectionString: connectionString,
		DatabaseName:     databaseName,
		CollectionName:   "unspent_outputs",
	}
}
