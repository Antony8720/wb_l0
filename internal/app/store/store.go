package store

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	config *Config
	db     *sql.DB
	mutex  sync.RWMutex
	Cache  *Cache
}

func New(config *Config) *Store {
	return &Store{
		config: config,
		Cache:  NewCache(),
	}
}

func (s *Store) Open() error {
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

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) GetOrdersAll() error {
	rows, err := s.db.Query("SELECT * FROM orders")

	defer rows.Close()
	if err != nil {
		log.Println("rows error")
		return err
	}

	for rows.Next() {
		var id string
		var data []byte
		err = rows.Scan(&id, &data)
		if err != nil {
			log.Println("rows scan error")
			return err
		}

		s.Cache.Set(id, data)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetOrderByID(id string) ([]byte, bool) {
	res, ok := s.Cache.Get(id)
	if !ok {
		log.Println("get !ok")
		return nil, false
	}
	return res, true
}

func (s *Store) AddOrder(id string, order []byte) error {
	if _, ok := s.Cache.data[id]; ok {
		return fmt.Errorf("Error: This ID is already use")
	} else {
		s.Cache.data[id] = order
	}
	query := "INSERT INTO orders (uid, data) VALUES ($1, $2)"
	_, err := s.db.Exec(query, id, order)
	if err != nil {
		return fmt.Errorf("Error: Inserting error: %v", err)
	}

	return nil
}
