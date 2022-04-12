package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
)

var db *sql.DB
var errj error

func main() {
	db, errj = sql.Open("mysql", "Z85D787U9Y:muEC8eUbmd@tcp(remotemysql.com:3306)/Z85D787U9Y")

	if errj != nil {
		panic(errj.Error())
	}

	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/calendario/event", addEvent).Methods("POST")
	http.ListenAndServe(":8000", router)
}

func addEvent(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	id_uc := keyVal["id_uc"]
	date := keyVal["date"]
	name := keyVal["name"]
	description := keyVal["description"]

	log.Println(id_uc)
	log.Println(date)
	log.Println(name)
	log.Println(description)
	stmt, err := db.Prepare("INSERT INTO calendario (id_uc, date, name, description) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(id_uc, date, name, description)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New post was created")
}

func getEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//RECEBE O JSON QUE VEM DA ROTA
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//COLOCA O JSON NUM MAP E ASSOCIA A VARIAVEIS
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	date := keyVal["date"]
	description := keyVal["description"]
	name := keyVal["name"]
	id_uc := keyVal["id_uc"]

	//EXECUTA A QUERY PARA O DB
	stmt := db.Query("SELECT * FROM calendario WHERE date=?", date)
	var dateComp date
	var descriptionComp string
	var nameComp string
	var id_ucComp int

	switch err := stmt.Scan(&dateComp); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(dateComp)
	default:
		panic(err)
	}
}
