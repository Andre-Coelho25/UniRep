package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id       int
	Name     string
	email    string
	password string
}

func main() {
	fmt.Println("Go MySQL Tutorial")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "Z85D787U9Y:muEC8eUbmd@tcp(remotemysql.com:3306)/Z85D787U9Y")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	results, err := db.Query("SELECT * FROM user")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var tag User
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.Id, &tag.Name, &tag.email, &tag.password)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(tag.Name)
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	fmt.Println("Olha que coisa linda")

}
