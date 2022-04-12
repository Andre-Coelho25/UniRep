package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"

	gomail "gopkg.in/gomail.v2"
)

var htmlBody = `
<html>
<head>
   <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
   <title>Hello, World</title>
</head>
<body>
   <p>Clique neste link1 para alterar a password: http://127.0.0.1:5500/FrontEnd/password.html</p>
   <img src="https://atlas-content-cdn.pixelsquid.com/stock-images/stickman-k1m1Ja3-600.jpg" width="400px" height="365px">
</body>
`

type User struct {
	Name     string `json:"name", db:"name"`
	email    string `json:"email", db:"email"`
	password string `json:"password", db:"password"`
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
	router.HandleFunc("/user/register", Signup)
	router.HandleFunc("/user/login", Signin)
	router.HandleFunc("/user/recover", recover)
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

func Signup(w http.ResponseWriter, r *http.Request) {

	//RECEBE O JSON QUE VEM DA ROTA
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//COLOCA O JSON NUM MAP E ASSOCIA A VARIAVEIS
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	email := keyVal["email"]
	password := keyVal["password"]

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	log.Println(name)
	log.Println(email)
	log.Println(hashedPassword)

	//PREPARA A QUERY PARA O DB
	stmt, err := db.Prepare("INSERT INTO user (name, email, password) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	//EXECUTA A QUERY DO DB COM OS PARAMETROS SENDO AS VARIAVEIS ASSOCIADAS PREVIAMENTE
	_, err = stmt.Exec(name, email, hashedPassword)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New post was created")
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}

func Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//RECEBE O JSON QUE VEM DA ROTA
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//COLOCA O JSON NUM MAP E ASSOCIA A VARIAVEIS
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	password := keyVal["password"]
	email := keyVal["email"]

	log.Println(email)

	//PREPARA A QUERY PARA O DB
	stmt := db.QueryRow("SELECT password FROM user WHERE email=?", email)
	var passwordComp string
	switch err := stmt.Scan(&passwordComp); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(passwordComp)
	default:
		panic(err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(passwordComp), []byte(password)); err != nil {
		// If the two passwords don't match, return a 401 status
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func recover(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//RECEBE O JSON QUE VEM DA ROTA
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//COLOCA O JSON NUM MAP E ASSOCIA A VARIAVEIS
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	email := keyVal["email"]

	msg := gomail.NewMessage()
	msg.SetHeader("From", "Suporte.UniRepository@gmail.com")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Teste")
	msg.Embed("filipeok.png")
	//	msg.SetBody("text/html", "<b>Clique neste link para alterar a password: http://127.0.0.1:5500/FrontEnd/password.html</b> <img source="FrontEnd\imgs\filipeok.png">")
	msg.SetBody("text/html", "<b>Clique neste link1 para alterar a password: http://127.0.0.1:5500/FrontEnd/password.html</b> <br><br><br> <img src='cid:filipeok.png' width='400px' height='365px'>")

	n := gomail.NewPlainDialer("smtp.gmail.com", 587, "Suporte.UniRepository@gmail.com", "OFilipeEBonito")

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}
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
