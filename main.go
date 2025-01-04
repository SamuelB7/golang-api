package main

import (
	"api/src/config"
	"api/src/database"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.LoadEnvs()
	database.Connect()
	r := router.GenerateRouter()
	fmt.Println("API running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
