package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
)

var dbj *sql.DB
var errj error

func mainj() {
	fmt.Println("Go MySQL Tutorial")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	dbj, errj = sql.Open("mysql", "Z85D787U9Y:muEC8eUbmd@tcp(remotemysql.com:3306)/Z85D787U9Y")

	// if there is an error opening the connection, handle it
	if errj != nil {
		panic(errj.Error())
	}

	defer dbj.Close()

	router := mux.NewRouter()
	router.HandleFunc("/posts", getUser).Methods("GET")
	router.HandleFunc("/user/{id}", getUserId).Methods("GET")

	http.ListenAndServe(":8000", router)

	// defer the close till after the main function has finished
	// executing

	fmt.Println("Olha que coisa linda")

}