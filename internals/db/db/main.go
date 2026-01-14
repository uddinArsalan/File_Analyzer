package db

import (
	"database/sql"
	"file-analyzer/internals/domain"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DBClient struct {
	l  *log.Logger
	db *sql.DB
}

func NewDBConnection(l *log.Logger) (*DBClient, error) {
	connStr := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Error Initialising Connection %v", err.Error())
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error Ping Connection %v", err.Error())
	}
	return &DBClient{l, db}, nil
}

func (dbClient *DBClient) CloseDB() error {
	return dbClient.db.Close()
}

func (dbClient *DBClient) FindUserByEmail(email string) (domain.User, error) {
	var user domain.User
	query := `SELECT user_id,name,email, password_hash
    		FROM users
    		WHERE email = $1`
	err := dbClient.db.
		QueryRow(query, email).
		Scan(&user.UserID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (dbClient *DBClient) FindUserById(userId string) (domain.User, error) {
	var user domain.User
	query := `
    		SELECT user_id,name,email, password_hash
    		FROM users
    		WHERE user_id = $1
		`
	err := dbClient.db.
		QueryRow(query, userId).
		Scan(&user.UserID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (dbClient *DBClient) FindUserByToken(tokenHash string) (domain.RefreshToken, error) {
	var token domain.RefreshToken
	query := `
    		SELECT id,user_id,expires_at,revoked_at
    		FROM refresh_tokens
    		WHERE token_hash = $1
		`
	err := dbClient.db.
		QueryRow(query, tokenHash).
		Scan(token.ID, &token.UserID, &token.ExpiresAt, &token.RevokedAt)
	if err != nil {
		return domain.RefreshToken{}, err
	}
	return token, nil
}

func (dbClient *DBClient) InsertUser(user domain.User) error {
	query := `
    		INSERT INTO users (name,email,password_hash) VALUES ($1,$2,$3)
		`
	_, err := dbClient.db.Exec(query, user.Name, user.Email, user.PasswordHash)
	return err
}

func (dbClient *DBClient) InsertRefreshToken(tokenHash string, userID int64, ttl time.Duration) (int64, error) {
	query := `
	INSERT INTO refresh_tokens (user_id,token_hash,expires_at) VALUES($1,$2,$3) RETURNING id
	`
	var id int64
	err := dbClient.db.QueryRow(query, userID, tokenHash, time.Now().Add(ttl)).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (dbClient *DBClient) RevokeRefreshToken(oldTokenID int64, newTokenID int64) error {
	query := `
	UPDATE refresh_tokens SET revoked_at = NOW(),replaced_by_token_id = $1 WHERE id = $2
	`
	_, err := dbClient.db.Exec(query, newTokenID, oldTokenID)
	return err
}
