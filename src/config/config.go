package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DatabaseUrl      = ""
	DatabasePort     = ""
	DatabaseHost     = ""
	DatabaseName     = ""
	DatabaseUser     = ""
	DatabasePassword = ""
)

func LoadEnvs() {
	var error error

	if error = godotenv.Load(); error != nil {
		log.Fatal(error)
	}

	DatabasePort = os.Getenv("DATABASE_PORT")
	DatabaseHost = os.Getenv("DATABASE_HOST")
	DatabaseName = os.Getenv("DATABASE_NAME")
	DatabaseUser = os.Getenv("DATABASE_USER")
	DatabasePassword = os.Getenv("DATABASE_PASSWORD")
}
