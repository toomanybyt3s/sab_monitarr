package main

import (
	"log"

	"github.com/toomanybyt3s/sab_monitarr/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
