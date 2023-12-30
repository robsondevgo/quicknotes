package main

import (
	"fmt"
	"html/template"
	"os"
)

type TemplateData struct {
	Nome string
	Age  int
}

func main() {
	t, err := template.ParseFiles("layout1.html", "footer.html", "header.html")

	fmt.Println(t.Name())

	fmt.Println(t.DefinedTemplates())

	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(os.Stdout, "layout1.html", "2023")
	if err != nil {
		panic(err)
	}
}
