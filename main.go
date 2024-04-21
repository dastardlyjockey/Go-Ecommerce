package main

import (
	"github.com/dastardlyjockey/ecommerce/database"
	"github.com/dastardlyjockey/ecommerce/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}

	port := os.Getenv("PORT")

	server := gin.New()

	server.Use(gin.Logger())

	routes.Routes(server)

	log.Println("Starting the server on port: ", port)
	err = server.Run(":" + port)
	if err != nil {
		log.Fatal("Server error: ", err)
	}

	defer database.CloseClient(database.Client)

}
