package ledger

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type LedgerDb struct {
	db *sql.DB

	LedgerTable *LedgerTable
}

func NewLedgerDb(dbPath string) (*LedgerDb, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	ledgerDb := &LedgerDb{
		db:          db,
		LedgerTable: NewLedgerTable(db),
	}

	err = ledgerDb.applyMigrations()
	if err != nil {
		return nil, err
	}

	return ledgerDb, nil
}

func (l *LedgerDb) Disconnect() {
	err := l.db.Close()
	if err != nil {
		panic(err)
	}
}
