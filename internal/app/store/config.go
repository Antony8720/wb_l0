package store

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
}

func NewConfig() *Config{
	cfg := Config{}
	flag.StringVar(&cfg.DatabaseURL, "b", "", "Database URL")
	flag.Parse()
	cfg.chooseDBURL()
	return &cfg
}

func (config *Config) chooseDBURL() {
	if config.DatabaseURL != "" {
		return
	}
	url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("Database URL not found")
		return
	}
	config.DatabaseURL = url
}