package main

import "wb_l0/internal/app/apiserver"

func main() {
	config := apiserver.NewConfig()
	s := apiserver.New(config)
	
}
