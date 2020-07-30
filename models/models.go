package models

import (
	"crypto/rand"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type DB struct {
	*sqlx.DB
}

type ConfigDB struct {
	User     string `env:"POSTGRES_USER,required"`
	Pass     string `env:"POSTGRES_PASS,required"`
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	Database string `env:"POSTGRES_DATABASE,required"`
	SslMode  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
}

type User struct {
	ID        string
	Username  string
	CreatedAt time.Time `db:"created_at"`
}

type Chat struct {
	ID        uint64
	Name      string
	Users     []uint64
	CreatedAt time.Time `db:"created_at"`
}

type Message struct {
	ID        uint64
	Chat      uint64
	Author    string
	Text      string
	CreatedAt time.Time `db:"created_at"`
}

type Interactor interface {
	CreateUser(username string) (string, error)
	CreateChat(chatName string, userIDs []string) (uint64, error)
	CreateMessage(chatID uint64, authorID string, msg string) (uint64, error)
	GetUserChats(userID string) ([]Chat, error)
	GetChatMessages(chatID uint64) ([]Message, error)
}

func ParseConfigDB() (*ConfigDB, error) {
	cfg := &ConfigDB{}
	if err := env.Parse(cfg); err != nil {
		log.Println(err)
		return nil, err
	}
	return cfg, nil
}

func InitDB(config *ConfigDB) (*DB, error) {
	sqlxDB, err := sqlx.Connect("pgx", "host=localhost port=5432 user=allerria password=root dbname=messenger sslmode=disable")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db := &DB{sqlxDB}
	return db, nil
}

func generateID() (string, error) {
	var s string
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Error: can't generate id")
		return s, err
	}
	s = fmt.Sprintf("%x", b)
	return s, nil
}

func (db *DB) CreateUser(username string) (string, error) {
	id, err := generateID()
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, err = db.Exec("INSERT INTO users (id, username) VALUES ($1, $2) RETURNING id", id, username)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return id, nil
}

func (db *DB) CreateChat(chatName string, userIDs []string) (uint64, error) {

	return 0, nil
}

func (db *DB) CreateMessage(chatID uint64, authorID string, msg string) (uint64, error) {
	return 0, nil
}

func (db *DB) GetUserChats(userID string) ([]Chat, error) {
	return nil, nil
}

func (db *DB) GetChatMessages(chatID uint64) ([]Message, error) {
	return nil, nil
}
