package main

import (
	"os"
	"text/template"
)

type Course struct {
	Name string
	Hour int
}

func main() {
	course := Course{"Go", 40}
	t := template.Must(template.New("CourseTemp2").Parse("Course: {{.Name}} - Hour: {{.Hour}}"))
	err := t.Execute(os.Stdout, course)
	if err != nil {
		panic(err)
	}
}
