package apiserver

import(
	"flag"
	"os"
)

type Config struct {
	BindAddress string
}

func NewConfig() *Config {
	cfg := Config{}
	flag.StringVar(&cfg.BindAddress, "a", "", "Bind address of the HTTP server")
	flag.Parse()
	cfg.chooseAddress()
	return &cfg
}

func (config *Config) chooseAddress() {
	if config.BindAddress != "" {
		return
	}
	address, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		address = ":8080"
	}
	config.BindAddress = address
}