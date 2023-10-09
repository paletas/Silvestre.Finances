package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
)

type dbVersion struct {
	Version    int
	UpScript   string
	DownScript string
}

func (version dbVersion) applyUpScript(db *sql.Tx) error {
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

func (version dbVersion) applyDownScript(db *sql.Tx) error {
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

func (ledgerDb FinancesDb) applyMigrations() error {
	log.Println("Applying migrations...")

	migrations, err := setupMigrations()
	if err != nil {
		return err
	}

	err = ensureMigrationReady(ledgerDb.db)
	if err != nil {
		return err
	}

	log.Println("Checking database version...")

	version, err := getDatabaseVersion(ledgerDb.db)
	if err != nil {
		return err
	}

	log.Printf("Database version: %d\n", version)

	conn, err := ledgerDb.db.Conn(context.Background())
	if err != nil {
		return err
	}

	tx, err := conn.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for i := version + 1; ; i++ {
		dbVersion, ok := migrations[i]
		if !ok {
			break
		}

		log.Printf("Applying migration %d\n", dbVersion.Version)

		err = migrations[i].applyUpScript(tx)
		if err != nil {
			log.Printf("Error applying migration %d\n", dbVersion.Version)
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	log.Println("Migrations applied successfully")

	return nil
}

func ensureMigrationReady(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS migration_logs(version INTEGER PRIMARY KEY, date TEXT)")
	if err != nil {
		return err
	}
	return nil
}

func setupMigrations() (map[int]dbVersion, error) {
	migrations := make(map[int]dbVersion)

	version, err := setupVersion(1)
	if err != nil {
		return nil, err
	}

	migrations[1] = version

	return migrations, nil
}

func getDatabaseVersion(db *sql.DB) (int, error) {
	var version int
	row := db.QueryRow("SELECT version FROM migration_logs ORDER BY version DESC LIMIT 1")
	err := row.Scan(&version)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return 0, err
		}
	}

	return version, nil
}

func setupVersion(version int) (dbVersion, error) {
	upScript, err := migrationScripts.ReadFile(fmt.Sprintf("migrations/version%d_up.sql", version))
	if err != nil {
		return dbVersion{}, err
	}

	downScript, err := migrationScripts.ReadFile(fmt.Sprintf("migrations/version%d_down.sql", version))
	if err != nil {
		return dbVersion{}, err
	}

	return dbVersion{
		Version:    version,
		UpScript:   string(upScript),
		DownScript: string(downScript),
	}, nil
}
