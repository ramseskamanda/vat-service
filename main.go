package main

import (
	"log"

	"github.com/ramseskamanda/vat-service/server"
)

func main() {
	err := server.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}
