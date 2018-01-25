package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {


	connStr := "postgresql://localhost?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	

	// _, err = db.Query("CREATE TABLE users (id integer, name varchar(40));")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	_, err = db.Query("INSERT INTO users VALUES(1, 'samuel');")
	if err != nil {
		log.Fatal(err)
	}
	
	insertId := 1
	rows, err := db.Query("SELECT name FROM users WHERE id = $1", insertId)
	if err != nil {
		log.Fatal(err)
	}

	var name string

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(name)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}