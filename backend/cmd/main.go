package main

import (
	"log"
	"github.com/KuRoute/kuroute/backend/package/database"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("Starting server on :8080")
}