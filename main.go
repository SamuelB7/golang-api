package main

import (
	"api/src/config"
	"api/src/database"
	"api/src/router"
	"fmt"
	"log"
	"net/http"

	_ "api/docs"
)

// @title          DevBook API
// @version        1.0
// @description    A simple social media API
// @termsOfService http://swagger.io/terms/

// @contact.name  Samuel Belo
// @contact.email belo.samuel@gmail.com

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer JWT token authentication, type "Bearer {token}"
// @tokenUrl /login
func main() {
	config.LoadEnvs()
	database.Connect()
	r := router.GenerateRouter()
	fmt.Println("API running on port 8080 with base path /api")
	log.Fatal(http.ListenAndServe(":8080", r))
}
