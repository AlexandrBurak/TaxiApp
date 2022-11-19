package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upUsers, downUsers)
}

func upUsers(tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS users(
	"id" SERIAL PRIMARY KEY,
	"name" VARCHAR(50),
	"phone" VARCHAR(12),
	"email" VARCHAR(30),
	"password" VARCHAR(100)
	);`
	_, err := tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func downUsers(tx *sql.Tx) error {
	query := `DROP TABLE users;`
	_, err := tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
