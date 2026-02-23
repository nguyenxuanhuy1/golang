package main

import (
	"log"
	"os"

	"traingolang/internal/api/router"
	"traingolang/internal/config"
)

func main() {
	config.ConnectDB()
	r := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Println("Server is running on :" + port)
	r.Run(":" + port)
}
