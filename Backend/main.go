package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
)

type User struct {
	Id       int
	Name     string
	email    string
	password string
}

var db *sql.DB
var err error

func main() {
	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err = sql.Open("mysql", "Z85D787U9Y:muEC8eUbmd@tcp(remotemysql.com:3306)/Z85D787U9Y")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/posts", getUser).Methods("GET")
	router.HandleFunc("/user/{id}", getUserId).Methods("GET")

	http.ListenAndServe(":8000", router)

	// defer the close till after the main function has finished
	// executing

	fmt.Println("Olha que coisa linda")

}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tags []User

	results, err := db.Query("SELECT * FROM user")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer results.Close()

	for results.Next() {
		var tag User
		// for each row, scan the result into our tag composite object
		err := results.Scan(&tag.Id, &tag.Name, &tag.email, &tag.password)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(tag.Name)
		tags = append(tags, tag)
	}

	json.NewEncoder(w).Encode(tags)
}

func getUserId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var tags []User

	results, err := db.Query("SELECT * FROM user WHERE id_user=?", params["id"])
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer results.Close()

	for results.Next() {
		var tag User
		// for each row, scan the result into our tag composite object
		err := results.Scan(&tag.Id, &tag.Name, &tag.email, &tag.password)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		tags = append(tags, tag)
		log.Printf(tags[0].email)

	}

	json.NewEncoder(w).Encode(tags)
}
