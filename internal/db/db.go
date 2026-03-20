package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"io"
	"os"

	_ "modernc.org/sqlite"
)

// InitDatabase - Initing database by configs.
// Sets configs, make migrations and any things. Prepare and returs database
func InitDatabase(configs *config.DbCfg) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "app.db")
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	err = execMigration(db, configs.DbMigrationPath)
	if err != nil {
		return nil, fmt.Errorf("ExecMigration: %w", err)
	}

	return db, err
}
func execMigration(db *sql.DB, mgPath string) error {
	f, err := os.OpenFile(mgPath, os.O_RDONLY, 0755)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}
	defer f.Close()

	migrationData, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migrationData))
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
