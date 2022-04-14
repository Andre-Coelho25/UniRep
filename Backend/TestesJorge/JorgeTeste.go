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

//É sempre com letra maiúscula a iniciar no JSon
//P.S: Não sei porquê
type Event struct {
	Date        string `json:"date"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

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
	router.HandleFunc("/calendario/eventbydate/{date}", getEventByDate).Methods("GET")

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

func getEventByDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	//EXECUTA A QUERY PARA O DB
	rows, err := db.Query("SELECT date, description, name FROM calendario WHERE date=?", params["date"])
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Date, &event.Description, &event.Name)
		if err != nil {
			panic(err.Error())
		}
		events = append(events, event)
		log.Printf(events[0].Date)
		log.Printf(events[0].Description)
		log.Printf(events[0].Name)
	}
	json.NewEncoder(w).Encode(events)
}
