package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"strings"
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
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Chat struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Users     []string  `json:"users"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Message struct {
	ID        uint64    `json:"id"`
	Chat      uint64    `json:"chat"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ChatUsers struct {
	ChatID uint64 `db:"chat_id"`
	UserID string `db:"user_id"`
}

type Model interface {
	CreateUser(username string) (string, error)
	CreateChat(chatName string, userIDs []string) (uint64, error)
	CreateMessage(chatID uint64, authorID string, text string) (uint64, error)
	GetUserChats(userID string) ([]Chat, error)
	GetChatMessages(chatID uint64) ([]Message, error)
}

func ParseConfigDB() (*ConfigDB, error) {
	cfg := &ConfigDB{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func InitDB(config *ConfigDB) (*DB, error) {
	sqlxDB, err := sqlx.Connect("pgx", "host=localhost port=5432 user=allerria password=root dbname=messenger sslmode=disable")
	if err != nil {
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
		return s, err
	}
	s = fmt.Sprintf("%x", b)
	return s, nil
}

func (db *DB) CreateUser(username string) (string, error) {
	id, err := generateID()
	if err != nil {
		return "", err
	}
	_, err = db.Exec("INSERT INTO users (id, username) VALUES ($1, $2)", id, username)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (db *DB) CreateChat(chatName string, userIDs []string) (uint64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	rows, err := tx.Query("INSERT INTO chats(name) VALUES($1) RETURNING id", chatName)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return 0, err
	}
	if !rows.Next() {
		tx.Rollback()
		return 0, errors.New("db didn't return created chat id")
	}
	var id uint64
	err = rows.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	err = rows.Close()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	queryStr := "INSERT INTO chats_users (chat_id, user_id) VALUES "
	vals := []interface{}{}
	argCount := 1
	for _, user := range userIDs {
		queryStr += fmt.Sprintf("($%s, $%s),", strconv.Itoa(argCount), strconv.Itoa(argCount+1))
		argCount += 2
		vals = append(vals, id, user)
	}
	queryStr = strings.TrimSuffix(queryStr, ",")
	stmt, err := tx.Prepare(queryStr)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	_, err = stmt.Exec(vals...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, nil
}

func (db *DB) CreateMessage(chatID uint64, authorID string, text string) (uint64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	rows, err := tx.Query("INSERT INTO messages (chat, author, text) VALUES($1, $2, $3) RETURNING id", chatID, authorID, text)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if !rows.Next() {
		tx.Rollback()
		return 0, errors.New("db didn't return created msg id")
	}
	var id uint64
	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *DB) GetUserChats(userID string) ([]Chat, error) {
	chats := []Chat{}
	chatUsers := []ChatUsers{}
	c := make(map[uint64]*Chat)
	err := db.Select(&chats, `SELECT id, name, created_at
FROM (SELECT id,
             name,
             created_at,
             (SELECT MAX(created_at) OVER (PARTITION BY id) FROM messages WHERE chat = id) AS last_msg_time
      FROM chats
      WHERE id IN (SELECT chat_id FROM chats_users WHERE user_id = $1)
      ORDER BY last_msg_time DESC) as t`, userID)
	if err != nil {
		return chats, err
	}
	for i, chat := range chats {
		c[chat.ID] = &chats[i]
	}
	err = db.Select(&chatUsers, "SELECT chat_id, user_id FROM chats_users WHERE chat_id IN (SELECT chat_id FROM chats_users WHERE user_id = $1)", userID)
	if err != nil {
		return chats, err
	}
	for _, cu := range chatUsers {
		c[cu.ChatID].Users = append(c[cu.ChatID].Users, cu.UserID)
	}
	return chats, nil
}

func (db *DB) GetChatMessages(chatID uint64) ([]Message, error) {
	msgs := []Message{}
	err := db.Select(&msgs, "SELECT * FROM messages WHERE chat = $1", chatID)
	if err != nil {
		return msgs, err
	}
	return msgs, nil
}
