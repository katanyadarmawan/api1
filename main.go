package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

type User struct {
	ID    int    `form:"id" json:"id"`
	Name  string `form:"name" json:"name"`
	Email string `form:"email" json:"email"`
}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {

	db, err = gorm.Open("mysql", "root:@/api1?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("Failed to create a connection to database")
	}
	//nanti kita isi modelnya di sini

	db.AutoMigrate(&User{})
	handlerRequest()
}

func handlerRequest() {

	fmt.Println("Start develop at server http://localhost:9000")
	fmt.Println("Running")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/user", getUsers).Methods("GET")
	myRouter.HandleFunc("/user/{id}", getUser).Methods("GET")
	myRouter.HandleFunc("/user", createUser).Methods("POST")
	myRouter.HandleFunc("/user/{id}", updateUser).Methods("POST")
	myRouter.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9000", myRouter))

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hellow World!")
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	user := []User{}

	json.Marshal(&user)

	db.Find(&user)

	res := Response{Code: 200, Data: user, Message: "Get all data user success"}

	result, err := json.Marshal(res)

	if err != nil {
		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userID := vars["id"]

	var user User

	db.First(&user, userID)

	res := Response{Code: 200, Data: user, Message: "Get data user success"}

	result, err := json.Marshal(res)

	if err != nil {
		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func createUser(w http.ResponseWriter, r *http.Request) {

	data, _ := ioutil.ReadAll(r.Body)

	var user User
	json.Unmarshal(data, &user)

	db.Create(&user)

	res := Response{Code: 200, Data: user, Message: "Create data user success"}

	result, err := json.Marshal(res)

	if err != nil {
		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userID := vars["id"]

	data, _ := ioutil.ReadAll(r.Body)

	var newData User
	json.Unmarshal(data, &newData)

	var user User

	db.First(&user, userID)
	db.Model(&user).Update(newData)

	res := Response{Code: 200, Data: user, Message: "Update data user success"}

	result, err := json.Marshal(res)

	if err != nil {
		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userID := vars["id"]

	var user User

	db.Delete(&user, userID)

	res := Response{Code: 200, Message: "Delete data user success"}

	result, err := json.Marshal(res)

	if err != nil {
		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
