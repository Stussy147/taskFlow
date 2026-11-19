package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "host=localhost port=5432 user=postgres password=2512 dbname=taskFlow sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка инициализации БД: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}
	fmt.Println("Успешное подключение к PostgreSQL!")
}
