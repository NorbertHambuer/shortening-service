package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/norby7/shortening-service/entities"
	"io/ioutil"
	"os"
)

type SqliteStorage struct {
	Handler *sql.DB
}

var (
	SqlOpen = sql.Open
)

// NewSqliteStorage connects to a sqlite database and returns a repository object that contains the database connection handler
func NewSqliteStorage(p string, maxConns int) (*SqliteStorage, error) {
	db, err := SqlOpen("sqlite3", p)
	if err != nil {
		return nil, fmt.Errorf("unable to open sqlite database: %s", err.Error())
	}

	if maxConns != 0 {
		db.SetMaxIdleConns(maxConns)
		db.SetMaxOpenConns(maxConns)
	}

	return &SqliteStorage{Handler: db}, nil
}

// CreateDatabase checks if the database file exists and creates one if it doesn't
func CreateDatabase(p string) error {
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		if _, err = os.Create(p); err != nil {
			return fmt.Errorf("unable to create sqlite database file: %s", err.Error())
		}
	}

	return nil
}

// ValidateSchema checks if the urls schemas exist and creates them if they don't exist
func ValidateSchema(db *sql.DB) error {
	dbExists, err := schemaExists(db)
	if err != nil {
		return fmt.Errorf("unable to check if database schema exists: %s", err.Error())
	}

	if !dbExists {
		err = initializeSchema(db)
		if err != nil {
			return fmt.Errorf("unable to create database schema: %s", err.Error())
		}
	}

	return nil
}

// schemaExists checks if the urls tables exist
func schemaExists(handler *sql.DB) (bool, error) {
	var n string
	if err := handler.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='urls';`).Scan(&n); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// initializeSchema reads the schema.sql file and executes the queries inside it
func initializeSchema(handler *sql.DB) error {
	c, err := ioutil.ReadFile("./database/sqlite/schema.sql")
	if err != nil {
		return fmt.Errorf("unable to open databse schema sql script: %s", err.Error())
	}

	sqlScript := string(c)

	// begin transaction
	tx, err := handler.Begin()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to start transaction: %s", err.Error())
	}

	// execute insert urls statement
	_, err = tx.Exec(sqlScript)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to execute create database queries: %s", err.Error())
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unable to commit transation: %s", err.Error())
	}

	return nil
}

// Add inserts a new url into the database and returns an error in case something went wrong
func (s *SqliteStorage) Add(url *entities.Url) error {
	res, err := s.Handler.Exec(`INSERT INTO urls (code, url, counter, shortUrl, domain) VALUES (?, ?, ?, ?, ?)`, url.Code, url.Url, url.Counter, url.ShortUrl, url.Domain)
	if err != nil {
		return err
	}

	// get new url id
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("unable to get last inserted id: %s", err.Error())
	}

	// set the Url new Id
	url.Id = id

	return nil
}

// Delete removes a url from the database based on the given Id
func (s *SqliteStorage) Delete(id int64) error {
	if _, err := s.Handler.Exec(`DELETE FROM urls WHERE id = ?`, id); err != nil {
		return err
	}

	return nil
}

// GetUrlByCode returns a url from the database with the given code
func (s *SqliteStorage) GetUrlByCode(code string) (string, error) {
	var url string
	if err := s.Handler.QueryRow(`SELECT url FROM urls WHERE code = ?`, code).Scan(&url); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	return url, nil
}

// GetById returns a url from the database with the given id
func (s *SqliteStorage) GetById(id int64) (entities.Url, error) {
	var u entities.Url
	if err := s.Handler.QueryRow(`SELECT * FROM urls WHERE id = ?`, id).Scan(&u.Id, &u.Code, &u.Url, &u.ShortUrl, &u.Domain, &u.Counter); err != nil {
		if err == sql.ErrNoRows {
			return entities.Url{}, nil
		}

		return entities.Url{}, err
	}

	return u, nil
}

// GetByUrl returns a url object from the database with the given url
func (s *SqliteStorage) GetByUrl(url string) (entities.Url, error) {
	var u entities.Url
	if err := s.Handler.QueryRow(`SELECT * FROM urls WHERE url = ?`, url).Scan(&u.Id, &u.Code, &u.Url, &u.ShortUrl, &u.Domain, &u.Counter); err != nil {
		if err == sql.ErrNoRows {
			return entities.Url{}, nil
		}

		return entities.Url{}, err
	}

	return u, nil
}

// IncrementCounter increments the counter for the given code
func (s *SqliteStorage) IncrementCounter(code string) error {
	if _, err := s.Handler.Exec(`UPDATE urls SET counter = counter + 1 WHERE code = ?`, code); err != nil {
		return err
	}

	return nil
}
