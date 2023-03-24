package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Start Request")
	defer log.Println("End Request")

	select {
	case <-time.After(5 * time.Second):
		log.Println("Success Request")
		w.Write([]byte("Success Request"))
	case <-ctx.Done():
		log.Println("Request canceled by client")
	}
}
