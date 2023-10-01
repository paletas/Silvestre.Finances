package sqlite

import (
	"database/sql"
	"embed"
)

type dbVersion struct {
	Version    int
	UpScript   byte[]
	DownScript byte[]
}

func (version *dbVersion) applyUpScript(db *sql.DB) {
	_, err := db.Exec(version.UpScript)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO migration_logs(version, date) VALUES(?, datetime('now'))", version.Version)
	if err != nil {
		return err
	}

	return nil
}

func (version *dbVersion) applyDownScript(db *sql.DB) {
	_, err := db.Exec(version.DownScript)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM migration_logs WHERE version = ?", version.Version)
	if err != nil {
		return err
	}

	return nil
}

var (
	//go:embed "migrations/*"
	migrationScripts embed.FS
)

func (ledgerDb SqliteLedger) applyMigrations() error {
	migrations := setupMigrations()

	err := ensureMigrationReady(ledgerDb.db)
	if err != nil {
		return err
	}

	version, err := getDatabaseVersion(ledgerDb.db)
	if err != nil {
		return err
	}

	for i := version; i < len(migrations); i++ {
		err = migrations[i].applyUpScript(ledgerDb.db)
		if err != nil {
			return err
		}
	}
}

func ensureMigrationReady(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS migration_logs(version INTEGER PRIMARY KEY, date TEXT)")
	if err != nil {
		return err
	}
	return nil
}

func setupMigrations() map[int]dbVersion {
	migrations := make(map[int]dbVersion)

	migrations[1] = dbVersion{
		Version:    1,
		UpScript:   migrationScripts.ReadFile("migrations/ledger/version1_up.sql"),
		DownScript: migrationScripts.ReadFile("migrations/ledger/version1_down.sql"),
	}

	return migrations
}

func getDatabaseVersion(db *sql.DB) (int, error) {
	var version int
	row := db.QueryRow("SELECT version FROM migration_logs ORDER BY version DESC LIMIT 1")
	err := row.Scan(&version)
	if err != nil {
		return 0, err
	}

	return version, nil
}