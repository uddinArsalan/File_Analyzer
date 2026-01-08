package db

import (
	"database/sql"
	"file-analyzer/internals/domain"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DBClient struct {
	l  *log.Logger
	db *sql.DB
}

func NewDBConnection(l *log.Logger) (*DBClient, error) {
	connStr := os.Getenv("DB_CONN_STR")
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

func (dbClient *DBClient) InsertUser(user domain.User) error {
	query := `
    		INSERT INTO TABLE users (name,email,password_hash) VALUES ($1,$2,$3)
		`
	_, err := dbClient.db.Exec(query, user.Name, user.Email, user.PasswordHash)
	return err
}
