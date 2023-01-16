package apiserver

import (
	"flag"
	"os"
	"wb_l0/internal/app/store"
)

type Config struct {
	BindAddress string
	Store *store.Config
}

func NewConfig() *Config {
	cfg := Config{
		Store: store.NewConfig(),
	}
	flag.StringVar(&cfg.BindAddress, "a", "", "Bind address of the HTTP server")
	flag.Parse()
	cfg.chooseAddress()
	return &cfg
}

func (config *Config) chooseAddress() {
	if config.BindAddress != "" {
		return
	}
	address, ok := os.LookupEnv("BIND_ADDRESS")
	if !ok {
		address = ":8080"
	}
	config.BindAddress = address
}