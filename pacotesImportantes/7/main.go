package main

import "net/http"

type Bar struct{}

type Foo struct {
	fuzzy string
}

func (f Foo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(f.fuzzy))
}

func (b Bar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Mux says Its Not a Foo"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mux at home"))
	})
	mux.HandleFunc("/hum", HumHandler)
	mux.Handle("/bar", Bar{})
	mux.Handle("/foo", Foo{fuzzy: "Mux says Its Not a Bar"})
	http.ListenAndServe(":8080", mux)
}

func HumHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Mux says Hummmmmmmmmmmm"))
}
