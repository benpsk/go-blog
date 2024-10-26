package main

import (
	"context"
	"log"
	"net/http"

	"github.com/benpsk/go-blog/config"
	"github.com/benpsk/go-blog/database"
	"github.com/benpsk/go-blog/internal"
)

func main() {
	conn, err := database.Connect(config.DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
	r := internal.Router(conn)

	log.Printf("Server running on: %v", config.PORT)
	log.Fatal(http.ListenAndServe(":"+config.PORT, r))
}
