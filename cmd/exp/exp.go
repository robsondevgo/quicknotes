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
	t, err := template.ParseFiles("hello.html")

	fmt.Println(t.Name()) //hello.html

	if err != nil {
		panic(err)
	}
	data := TemplateData{Nome: "Robson"}
	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
