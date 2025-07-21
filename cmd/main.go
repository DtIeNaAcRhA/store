package main

import (
	"fmt"
	"log"
	"net/http"
	"store/internal/config"
	"store/internal/database"
	"store/internal/handler"
)

func main() {
	config.Load()

	err := database.InitDB(config.AppConfig.DBPath)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}

	defer database.DB.Close()

	router := handler.NewRouter()

	port := fmt.Sprintf(":%s", config.AppConfig.Port)

	log.Printf("Server started at %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
