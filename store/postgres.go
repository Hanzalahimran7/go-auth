package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/hanzalahimran7/go-auth/model"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(host, port, user, password, dbName string) *PostgresDB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to Database!")
	return &PostgresDB{db: db}
}

func (p *PostgresDB) RunMigration() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP
	);
	`
	_, err := p.db.Exec(string(query))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}
	log.Println("Migration executed successfully!")
	return nil
}

func (p *PostgresDB) Signup(ctx context.Context, user *model.User) error {
	query := `
        INSERT INTO users (id, first_name, last_name, email, password, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := p.db.ExecContext(ctx, query, user.Id, user.FirstName, user.LastName, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("Failed to Signup: %w", err)
	}
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

func (p *PostgresDB) FindByEmail(ctx context.Context, email string) error {
	var emailInDB string
	if err := p.db.QueryRow("SELECT email from users where email = $1", email).Scan(&emailInDB); err != nil {
		return err
	}
	if emailInDB == email {
		return fmt.Errorf("email already exists")
	}
	return nil
}
