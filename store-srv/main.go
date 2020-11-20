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

type CreateStores struct {
	StoreID     int    `json:"store_id"`
	StoreName   string `json:"store_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

func getstoreshandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT store_id,store_name, address, phone_number  from stores")
	if err != nil {
		log.Fatal(err)
	}
	var stores []CreateStores
	defer rows.Close()
	for rows.Next() {
		var store CreateStores
		err := rows.Scan(&store.StoreID, &store.StoreName, &store.Address, &store.PhoneNumber)
		if err != nil {
			log.Fatal(err)
		}
		stores = append(stores, store)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(stores)

	js, err := json.Marshal(stores)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(js)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createstorehandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	var store CreateStores
	json.Unmarshal(reqBody, &store)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%v", store)
	senddata(store)

}
func senddata(store CreateStores) {

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	result, err := db.Query("INSERT INTO stores(store_name,address,phone_number) VALUES(?,?,?)", store.StoreName, store.Address, store.PhoneNumber)
	if err != nil {
		panic(err.Error())

	}
	fmt.Println(result)
}

func main() {
	dbUser := os.Getenv("MYSQL_USERNAME")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOSTNAME")
	dbDatabase := os.Getenv("MYSQL_DATABASE")
	dbURL = dbUser + ":" + dbPassword + "@tcp(" + dbHost + ")/" + dbDatabase
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/getstores", getstoreshandler).Methods("GET")
	router.HandleFunc("/createstore", createstorehandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
