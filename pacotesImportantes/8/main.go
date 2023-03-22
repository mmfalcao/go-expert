package main

import (
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./8/public"))
	mux := http.NewServeMux()
	mux.Handle("/", fileServer)
	log.Fatal(http.ListenAndServe(":8081", mux))
}
