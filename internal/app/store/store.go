package store

import (
	"database/sql"
	"sync"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	config *Config
	db *sql.DB
	mutex sync.RWMutex
	cache *Cache
}

type DatabaseJSON struct {
	userID int
	json string
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func(s *Store) Open() error{
	db, err := sql.Open("pgx", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	
	s.db = db

	return nil
}

func(s *Store) Close() {
	s.db.Close()
}

func(s *Store) GetOrdersAll() error{
	rows, err := s.db.Query()
	defer rows.Close()
	if err != nil {
		return err
	}
	
	for rows.Next() {
		var dbjs DatabaseJSON
		err = rows.Scan(&dbjs.userID, &dbjs.json)
		if err != nil {
			return err
		}

		s.cache.Set(dbjs.userID, dbjs.json)
	}

	err = rows.Err()
	if err != nil {
		return err
	}
	
	return nil
}

func(s *Store) GetOrderByID(int) (string, error) {

}