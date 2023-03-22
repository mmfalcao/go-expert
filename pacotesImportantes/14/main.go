package main

import (
	"html/template"
	"net/http"
	"strings"
)

type Course struct {
	Name string
	Hour int
}

type Courses []Course

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func main() {

	templates := []string{
		"header.html",
		"content.html",
		"footer.html",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.New("content.html")
		t.Funcs(template.FuncMap{"ToUpper": ToUpper})
		t = template.Must(t.ParseFiles(templates...))
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
