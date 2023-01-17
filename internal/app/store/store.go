package store

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	config *Config
	db *sql.DB
	mutex sync.RWMutex
	cache *Cache
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
	rows, err := s.db.Query("SELECT * FROM orders")
	defer rows.Close()
	if err != nil {
		return err
	}
	
	for rows.Next() {
		var id string
		var data []byte
		err = rows.Scan(&id, &data)
		if err != nil {
			return err
		}

		s.cache.Set(id, data)
	}

	err = rows.Err()
	if err != nil {
		return err
	}
	
	return nil
}

func(s *Store) GetOrderByID(id string) ([]byte, bool) {
	res, ok := s.cache.Get(id)
	if !ok {
		return nil, false
	}
	return res, true
}

func(s *Store) AddOrder(id string, order []byte) error {
	if _, ok := s.cache.data[id]; ok {
		return fmt.Errorf("Error: This ID is already use")
	} else {
		s.cache.data[id] = order
	} 

	query := "INSERT INTO orders (uid, data) VALUES ($1, $2)"
	s.db.Exec(query, id, order)

	return nil
}