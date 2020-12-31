package main

import (
	"github.com/DebanjanBarman/todo/db"
	"github.com/DebanjanBarman/todo/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	//Read environment variables
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading environment variables")
	}
	PORT := os.Getenv("PORT")
	MongodbConnectionUri := os.Getenv("MONGODB_CONNECTION_URI")

	//Connect to database
	db.ConnectDB(MongodbConnectionUri)

	//Setup Routes
	router := routes.Routes()

	log.Printf("Listening at port %v ...\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
