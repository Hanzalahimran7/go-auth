package store

import (
	"context"
	"database/sql"

	"github.com/hanzalahimran7/go-auth/model"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(host, port, user, password, db, sslMode string) *PostgresDB {
	return &PostgresDB{db: nil}
}

func (p *PostgresDB) runMigration() {}

func (p *PostgresDB) Signup(ctx context.Context, user *model.User) error {
	return nil
}

func (p *PostgresDB) Login(ctx context.Context, username string, password string) (*model.User, error) {
	return nil, nil
}

func (p *PostgresDB) GetUser(ctx context.Context, username string) (*model.User, error) {
	return nil, nil
}

func (p *PostgresDB) GetUsers(ctx context.Context) ([]*model.User, error) {
	return nil, nil
}

func (p *PostgresDB) DeleteUser(ctx context.Context, username string) error {
	return nil
}

func (p *PostgresDB) UpdateUser(ctx context.Context, user *model.User) error {
	return nil
}
