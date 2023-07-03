package main

import(
	"net/http"
	"log"
)

func main(){

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting Server on Port :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}