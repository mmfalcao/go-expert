package main

import (
	"os"
	"text/template"
)

type Course struct {
	Name string
	Hour int
}

type Courses []Course

func main() {
	t := template.Must(template.New("template.html").ParseFiles("template.html"))
	err := t.Execute(os.Stdout, Courses{
		{"Go", 40},
		{"Java", 80},
		{"Python", 20},
	})
	if err != nil {
		panic(err)
	}
}
