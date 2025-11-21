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

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("db connection: ", err)
	}

	app := NewApp(db, conf.ApiPort)
	app.Run()

}
