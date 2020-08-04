package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
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

// sql Queries used by Model.
const (
	SqlQueryCreateUser      = "INSERT INTO users (id, username) VALUES ($1, $2)"
	SqlQueryCreateChat      = "INSERT INTO chats(name) VALUES($1) RETURNING id"
	SqlQueryInsertChatUsers = "INSERT INTO chats_users (chat_id, user_id) VALUES "
	SqlQueryCreateMessage   = "INSERT INTO messages (chat, author, text) VALUES($1, $2, $3) RETURNING id"
	SqlQuerySelectUserChats = `SELECT id, name, created_at
FROM (SELECT id,
             name,
             created_at,
             (SELECT MAX(created_at) OVER (PARTITION BY id) FROM messages WHERE chat = id) AS last_msg_time
      FROM chats
      WHERE id IN (SELECT chat_id FROM chats_users WHERE user_id = $1)
      ORDER BY last_msg_time DESC) as t`
	SqlQuerySelectChatsUsers   = "SELECT chat_id, user_id FROM chats_users WHERE chat_id IN (SELECT chat_id FROM chats_users WHERE user_id = $1)"
	SqlQuerySelectChatMessages = "SELECT * FROM messages WHERE chat = $1 ORDER BY created_at ASC"
	SqlQueryCheckIfUserExist   = "SELECT id FROM users WHERE id = $1"
	SqlQueryCheckIfChatExist   = "SELECT id FROM chats WHERE id = $1"
)

// Errors
var (
	// ErrUserNotExist is returned when user doesn't exist.
	ErrUserNotExist = errors.New("chat doesn't exist")

	// ErrChatNotExist is returned when chat doesn't exist.
	ErrChatNotExist = errors.New("user doesn't exist")
)

func ParseConfig() (*ConfigDB, error) {
	cfg := &ConfigDB{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func InitDB() (*DB, error) {
	cfg, err := ParseConfig()
	if err != nil {
		return nil, err
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Database, cfg.SslMode)
	sqlxDB, err := sqlx.Connect("pgx", connStr)
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
	if err := validateCreateUserInput(username); err != nil {
		return "", err
	}
	id, err := generateID()
	if err != nil {
		return "", err
	}
	_, err = db.Exec(SqlQueryCreateUser, id, username)
	if err != nil {
		return "", err
	}
	return id, nil
}

// createInsertChatUsersSqlQuery returns sqlQuery with multiple insert values and values to insert
func createInsertChatUsersSqlQuery(userIDs []string, chatID uint64) (string, []interface{}) {
	queryStr := SqlQueryInsertChatUsers
	values := []interface{}{}
	argCount := 1
	for _, user := range userIDs {
		queryStr += fmt.Sprintf("($%s, $%s),", strconv.Itoa(argCount), strconv.Itoa(argCount+1))
		argCount += 2
		values = append(values, chatID, user)
	}
	queryStr = strings.TrimSuffix(queryStr, ",")
	return queryStr, values
}

// CreateChat returns created chat ID if creates chat in database else returns error
func (db *DB) CreateChat(chatName string, userIDs []string) (uint64, error) {
	if err := validateCreateChatInput(chatName, userIDs); err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	var id uint64
	err = tx.QueryRow(SqlQueryCreateChat, chatName).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	queryStr, values := createInsertChatUsersSqlQuery(userIDs, id)
	stmt, err := tx.Prepare(queryStr)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	_, err = stmt.Exec(values...)
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
	if err := validateCreateMessageInput(chatID, authorID, text); err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	var id uint64
	err = tx.QueryRow(SqlQueryCreateMessage, chatID, authorID, text).Scan(&id)
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

func (db *DB) GetUserChats(userID string) ([]Chat, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	if exist, err := db.CheckUserExist(userID); err != nil {
		return nil, err
	} else if !exist {
		return nil, ErrUserNotExist
	}

	chats := []Chat{}
	chatUsers := []ChatUsers{}
	c := make(map[uint64]*Chat)

	err := db.Select(&chats, SqlQuerySelectUserChats, userID)
	if err != nil {
		return chats, err
	}

	for i, chat := range chats {
		c[chat.ID] = &chats[i]
	}
	err = db.Select(&chatUsers, SqlQuerySelectChatsUsers, userID)
	if err != nil {
		return chats, err
	}
	for _, cu := range chatUsers {
		c[cu.ChatID].Users = append(c[cu.ChatID].Users, cu.UserID)
	}
	return chats, nil
}

func (db *DB) GetChatMessages(chatID uint64) ([]Message, error) {
	if err := validateChatID(chatID); err != nil {
		return nil, err
	}

	if exist, err := db.CheckChatExist(chatID); err != nil {
		return nil, err
	} else if !exist {
		return nil, ErrChatNotExist
	}

	msgs := []Message{}
	err := db.Select(&msgs, SqlQuerySelectChatMessages, chatID)
	if err != nil {
		return msgs, err
	}
	return msgs, nil
}

func (db *DB) CheckUserExist(ID string) (bool, error) {
	newID := ""
	err := db.Select(&newID, SqlQueryCheckIfUserExist, ID)
	if err != nil {
		return false, err
	}
	return newID == ID, nil
}

func (db *DB) CheckChatExist(ID uint64) (bool, error) {
	var newID uint64
	err := db.Select(&newID, SqlQueryCheckIfChatExist, ID)
	if err != nil {
		return false, err
	}
	return newID == ID, nil
}
