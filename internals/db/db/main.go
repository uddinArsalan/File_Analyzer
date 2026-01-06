package db

import (
	"database/sql"
	"file-analyzer/models"
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
	return &DBClient{l, db}, nil
}

func (dbClient *DBClient) CloseDB() error {
	return dbClient.db.Close()
}

func (dbClient *DBClient) CheckUserExist(email string) (string, error) {
	var userId string
	query := "SELECT user_id FROM users WHERE email = $1"
	row := dbClient.db.QueryRow(query, email)
	err := row.Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (dbClient *DBClient) FindUserById(userId string) (models.User, error) {
	var user models.User
	query := `
    		SELECT user_id,name,email, password_hash
    		FROM users
    		WHERE user_id = $1
		`
	err := dbClient.db.
		QueryRow(query, userId).
		Scan(&user.UserId, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (dbClient *DBClient) InsertUser(user models.User) error {
	query := `
    		INSERT INTO TABLE users (name,email,password_hash) VALUES ($1,$2,$3)
		`
	_, err := dbClient.db.Exec(query, user.Name, user.Email, user.PasswordHash)
	return err
}
