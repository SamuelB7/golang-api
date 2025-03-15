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
	ApiSecret        = ""
)

func LoadEnvs() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	DatabasePort = os.Getenv("DATABASE_PORT")
	DatabaseHost = os.Getenv("DATABASE_HOST")
	DatabaseName = os.Getenv("DATABASE_NAME")
	DatabaseUser = os.Getenv("DATABASE_USER")
	DatabasePassword = os.Getenv("DATABASE_PASSWORD")
	ApiSecret = os.Getenv("API_SECRET")
}
