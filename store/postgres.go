package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hanzalahimran7/go-auth/model"
	"github.com/hanzalahimran7/go-auth/utils"
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
    CREATE TABLE IF NOT EXISTS refresh_tokens (
        id SERIAL PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users(id),
        token VARCHAR(255) NOT NULL,
        expires_at TIMESTAMP,
        revoked BOOLEAN DEFAULT FALSE,
        CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
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
		return fmt.Errorf("FAILED TO SIGN UP: %w", err)
	}
	return nil
}

func (p *PostgresDB) Login(ctx context.Context, email string, password string) (model.User, error) {
	user := model.User{}
	if err := p.db.QueryRow("SELECT * from users where email = $1", email).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	); err != nil {
		return model.User{}, err
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return model.User{}, fmt.Errorf("INVALID PASSWORD")
	}
	return user, nil
}

func (p *PostgresDB) GetUser(ctx context.Context, id string) (model.User, error) {
	user := model.User{}
	if err := p.db.QueryRow("SELECT id, first_name, last_name, email, created_at from users where id = $1", id).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
	); err != nil {
		return model.User{}, err
	}
	return user, nil
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

func (p *PostgresDB) CheckEmailExists(ctx context.Context, email string) error {
	var emailInDB string
	if err := p.db.QueryRow("SELECT email from users where email = $1", email).Scan(&emailInDB); err != nil {
		return err
	}
	if emailInDB == email {
		return fmt.Errorf("email already exists")
	}
	return nil
}

func (p *PostgresDB) StoreToken(ctx context.Context, jwt string, userId uuid.UUID, expiresAt time.Time) error {
	query := `
        INSERT INTO refresh_tokens (user_id, token, expires_at, revoked)
        VALUES ($1, $2, $3, $4)
    `
	_, err := p.db.ExecContext(ctx, query, userId, jwt, expiresAt, false)
	if err != nil {
		return fmt.Errorf("FAILED TO SAVE REFRESH TOKEN: %v", err)
	}
	return nil
}

func (p *PostgresDB) GetTokenFromDB(ctx context.Context, userId string) (string, error) {
	return "", nil
}
