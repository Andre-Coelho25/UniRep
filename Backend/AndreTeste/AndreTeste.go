package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
)

type Usera struct {
	Id       int
	Name     string
	email    string
	password string
}

type Imgpath struct {
	ID        string `json:"id"`
	ImagePath string `json:"imgpath"`
}

var db *sql.DB
var err error
var imgpaths []Imgpath

func main() {
	fmt.Println("Go MySQL Tutorial")

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
	router.HandleFunc("/user/name/{id}", changeName).Methods("PUT")
	router.HandleFunc("/user/img", uploadImg)
	http.ListenAndServe(":8000", router)

	// defer the close till after the main function has finished
	// executing

	fmt.Println("Olha que coisa linda")

}

func changeName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE user SET name = ? WHERE id_user = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)

	json.Unmarshal(body, &keyVal)

	newTitle := keyVal["name"]

	_, err = stmt.Exec(newTitle, params["id"])

	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Post with ID = %s was updated", params["id"])
}

//////////////////////////////////////////////FUTURO TODODODODODODO////////////////////////////////////////
func uploadImg(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(1 << 2)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	//dst, err := os.Create(handler.Filename)
	////////////////////////////////////////////////////
	dst, err := os.Create(filepath.Join("C:/Users/andre/Desktop/UniRep-1/Backend/AndreTeste/temp-img", filepath.Base(handler.Filename)))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		fmt.Println(err)
		return
	}
	////////////////////////////////////////////////////
	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a new buffer base on file size
	fInfo, _ := dst.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(dst)
	fReader.Read(buf)

	//convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	//fmt.Fprintf(w,imgBase64Str)
	fmt.Fprintf(w, imgBase64Str)

	//Decoding
	sDec, _ := base64.StdEncoding.DecodeString(imgBase64Str)
	fmt.Println(sDec)
	filepat := "\\ImgUploadApp\\temp-img\\" + handler.Filename
	fmt.Println(filepat)
	db, _ := sql.Open("mysql", "Z85D787U9Y:muEC8eUbmd@tcp(remotemysql.com:3306)/Z85D787U9Y")
	insert, err := db.Query("INSERT INTO img (filename,filepath,imgdata) VALUES (?,?,?) WHERE id_user = 1", handler.Filename, filepat, imgBase64Str)
	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

	fmt.Fprintf(w, "Successfully Uploaded File\n"+"")
}
