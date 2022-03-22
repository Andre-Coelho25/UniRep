package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
)

var dbj *sql.DB
var errj error

func main() {
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
	//router.HandleFunc("/posts", getUser).Methods("GET")
	//router.HandleFunc("/user/{id}", getUserId).Methods("GET")

	http.ListenAndServe(":8000", router)

	// defer the close till after the main function has finished
	// executing

	/*eventoID, err := adicionaEvento(Evento{
		Date:        time.Now(),
		Description: "Coisa Linda Zé Carlos",
		Name:        "Dia da Depressão",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Id do album: %v\n", eventoID)*/
}

/*type Evento struct {
	Date        time.Time
	Description string
	Id_data     int64
	Id_uc_cal   int64
	Name        string
}*/

/*func adicionaEvento(evento Evento) (int64, error) {
	result, err := dbj.Exec("INSERT INTO calendario (date, description, name) VALUES (?, ?, ?)", evento.Date, evento.Description, evento.Name)
	if err != nil {
		return 0, fmt.Errorf("addEvento: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("adicionaEvento: %v", err)
	}
	return id, nil
}*/

func adicionaEvento(w http.ResponseWriter, r *http.Request) {
	stmt, err := dbj.Prepare("INSERT INTO calendario(date, description, name) VALUES (?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["evento"]
	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}
