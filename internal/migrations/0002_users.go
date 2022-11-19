package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up, Down)
}

func Up(tx *sql.Tx) error {
	query := `ALTER TABLE users ALTER COLUMN password TYPE VARCHAR(100)`
	_, err := tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func Down(tx *sql.Tx) error {
	query := `ALTER TABLE users ALTER COLUMN password TYPE VARCHAR(50);`
	_, err := tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
