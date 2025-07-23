package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"final/pkg/app"
	"final/pkg/db"
)

func main() {
	app.Init()

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	err = db.Init()
	if err != nil {
		log.Printf("Failed to initialize DB: %v", err)
	}

	webDir := "web"

	value := os.Getenv("TODO_PORT")
	if value == "" {
		value = "7540"
	}

	addr := ":" + value

	fmt.Println("Server started on port :7540")

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
