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
	tmp := template.New("CourseTemplate")
	tmp, _ = tmp.Parse("Course: {{.Name}} - Hour: {{.Hour}}")
	err := tmp.Execute(os.Stdout, course)
	if err != nil {
		panic(err)
	}
}
