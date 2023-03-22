package main

import (
	"net/http"
	"text/template"
)

type Course struct {
	Name string
	Hour int
}

type Courses []Course

func main() {

	templates := []string{
		"header.html",
		"content.html",
		"footer.html",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("content.html").ParseFiles(templates...))
		err := t.Execute(w, Courses{
			{"Go", 40},
			{"Java", 80},
			{"Python", 20},
		})
		if err != nil {
			panic(err)
		}
	})
	http.ListenAndServe(":8282", nil)

}
