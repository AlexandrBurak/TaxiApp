package repository

import (
	"context"
	"database/sql"

	"github.com/AlexandrBurak/TaxiApp/internal/config"
	"github.com/AlexandrBurak/TaxiApp/internal/model"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(cfg config.DbConfig) (*Repository, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	err = makeMigrations(db)
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (repos *Repository) AddNewUser(ctx context.Context, user model.User) error {

	sqlStatement := `
INSERT INTO users (name, phone, email, password)
VALUES ($1, $2, $3, $4)`
	_, err := repos.db.ExecContext(ctx, sqlStatement, user.Name, user.Phone, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil

}

func (repos *Repository) GetUserPhoneAndPasswordByPhone(ctx context.Context, user model.User) (model.User, error) {
	var userDB model.User

	sqlStatement := `SELECT phone, password FROM users WHERE phone=$1;`
	row := repos.db.QueryRowContext(ctx, sqlStatement, user.Phone)

	err := row.Scan(&userDB.Phone, &userDB.Password)
	if err != nil {
		return model.User{}, err
	}

	return model.User{}, nil
}

func (repos *Repository) Exists(ctx context.Context, user model.User) (bool, error) {
	var userDB model.User
	sqlStatement := `SELECT phone, password FROM users WHERE phone=$1;`
	row := repos.db.QueryRowContext(ctx, sqlStatement, user.Phone)

	err := row.Scan(&userDB.Phone, &userDB.Password)
	if err != nil {
		return false, err
	}
	return true, nil
}

func makeMigrations(db *sql.DB) error {
	err := goose.Up(db, "./internal/migrations")
	if err != nil {
		return err
	}
	return nil
}
