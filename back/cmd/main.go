package main

import (
	"log"
	"back/config"
)

func main() {
	log.Println("Starting application...")
	
	if err := config.Init(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}