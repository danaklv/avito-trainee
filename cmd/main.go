package main

import (
	"database/sql"
	"log"
	config "pr-reviewer/configs"

	_ "github.com/lib/pq"
)

func main() {
	conf := config.Load()

	db, err := sql.Open("postgres", conf.DbConn)

	if err != nil {
		log.Fatal("db open: ", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println("failed to close db:", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatal("db connection: ", err)
	}

	app := NewApp(db, conf.ApiPort)
	app.Run()

}
