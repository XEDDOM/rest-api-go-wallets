package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"restapi/internal/api/router"
	"restapi/internal/repository/sqlconnect"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("File not found")
		return
	}
	_, err = sqlconnect.ConnectDb()
	if err != nil {
		log.Fatalln("Database connection error", err)
		return
	}
	port := ":" + os.Getenv("API_PORT")
	router := router.WalletRouter()
	fmt.Println("Server is running on port:", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}
}
