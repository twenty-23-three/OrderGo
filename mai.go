package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func mai() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS time (id INTEGER PRIMARY KEY, time DATETIME)")
	if err != nil {
		panic(err)
	}
	statement.Exec()

	statement, err = db.Prepare("INSERT INTO time (time) VALUES (?)")
	if err != nil {
		panic(err)
	}
	statement.Exec(time.Now().Add(time.Hour * 2))

	rows, _ := db.Query("SELECT id, time FROM time")
	var id int
	var cTime time.Time

	for rows.Next() {
		rows.Scan(&id, &cTime)
		fmt.Println(id, cTime)
	}
}
