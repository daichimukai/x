package main

import (
	"log"

	"github.com/daichimukai/x/syakyo/proglog/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("%w", err)
	}
}
