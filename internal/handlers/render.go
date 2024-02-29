package handlers

import (
	"bytes"
	"html/template"
	"net/http"
)

func render(w http.ResponseWriter, status int, page string, data any) error {
	files := []string{
		"views/templates/base.html",
	}
	files = append(files, "views/templates/pages/"+page)
	t, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}
	buff := &bytes.Buffer{}
	err = t.ExecuteTemplate(buff, "base", data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	buff.WriteTo(w)
	return nil
}
