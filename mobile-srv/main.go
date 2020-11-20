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

type Mobile struct {
	StoreID      int    `json:"store_id"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	CostPrice    int    `json:"cost_price"`
	SellingPrice int    `json:"selling_price"`
}

var (
	dbURL = "root:root@tcp(db:3306)/mobilestore"
)


func getmobilehandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT store_id,brand, model, costprice,sellingprice  from mobiles")
	if err != nil {
		log.Fatal(err)
	}
	var stores []Mobile
	defer rows.Close()
	for rows.Next() {
		var store Mobile
		err := rows.Scan(&store.StoreID, &store.Brand, &store.Model, &store.CostPrice, &store.SellingPrice)
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

func stores(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	var mobile Mobile
	json.Unmarshal(reqBody, &mobile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%v", mobile)
	senddata(mobile)

}
func senddata(mobile Mobile) {

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	result, err := db.Query("INSERT INTO mobiles(store_id,brand,model,costprice,sellingprice) VALUES(?,?,?,?,?)", mobile.StoreID, mobile.Brand, mobile.Model, mobile.CostPrice, mobile.SellingPrice)
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
	router.HandleFunc("/getmobile", getmobilehandler).Methods("GET")
	router.HandleFunc("/createmobile", stores).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
