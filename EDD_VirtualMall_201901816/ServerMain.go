package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	request()
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Servidor en Go")
}

func getArreglo(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "[1,2,3,4]")
}

func metodoPost(w http.ResponseWriter, r *http.Request){
	body, _ := ioutil.ReadAll(r.Body)
	var request string
	json.Unmarshal(body, &request)
	fmt.Println(request)
}

func request(){
	myrouter := mux.NewRouter().StrictSlash(true)
	myrouter.HandleFunc("/", homePage)
	myrouter.HandleFunc("/GetArreglo", getArreglo).Methods("GET")
	myrouter.HandleFunc("/MetodoPost", metodoPost).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", myrouter))
}
