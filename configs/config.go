package config

import (
	"fmt"
	"os"
)

type Conf struct {
	DbConn string
	ApiPort string
}

func Load() *Conf {

	dbuser := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	dbport := os.Getenv("DB_PORT")
	dbpass := os.Getenv("DB_PASS")
	dbhost := os.Getenv("DB_HOST")
	apiport := os.Getenv("API_PORT")

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport,dbuser,dbpass,dbname)

	return &Conf {
		DbConn: conn,
		ApiPort: apiport,
	}

}
