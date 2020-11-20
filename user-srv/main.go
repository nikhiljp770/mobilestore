package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	dbURL = "root:root@tcp(db:3306)/mobilestore"
)

type User struct {
	StoreID     int    `json:"store_id"`
	UserName   string  `json:"user_name"`
	Password     string `json:"password"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data ")
	}
	var userinput User
	err = json.Unmarshal(reqBody, &userinput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%v", userinput)
	senddata(userinput)
}
func senddata(userinput User) {

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Println(err.Error())
	}
	defer db.Close()
	result, err := db.Query("INSERT INTO users( store_id,user_name,password) VALUES(?,?,?)",userinput.StoreID,userinput.UserName, userinput.Password)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(result)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data ")
	}

	var inputuser User
    var dbuser User
	err = json.Unmarshal(reqBody, &inputuser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%v", inputuser)
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	row := db.QueryRow("SELECT store_id,user_name,password from users where user_name=?",inputuser.UserName)

	err = row.Scan(&dbuser.StoreID, &dbuser.UserName, &dbuser.Password)
		if err != nil {
			log.Fatal(err)
		}
		var responsemessage string
	if(dbuser.StoreID == inputuser.StoreID && dbuser.UserName == inputuser.UserName && dbuser.Password == inputuser.Password) {
		responsemessage = "login succesful"
	} else {
		responsemessage ="invalid user"

	}
	js, err := json.Marshal(responsemessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(js)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	dbUser := os.Getenv("MYSQL_USERNAME")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOSTNAME")
	dbDatabase := os.Getenv("MYSQL_DATABASE")
	dbURL = dbUser + ":" + dbPassword + "@tcp(" + dbHost + ")/" + dbDatabase
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/user/register", registerHandler).Methods("POST")
	router.HandleFunc("/user/login", loginHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
