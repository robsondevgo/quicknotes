package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/robsondevgo/quicknotes/internal/apperror"
)

type noteHandler struct{}

func NewNoteHandler() *noteHandler {
	return &noteHandler{}
}

func (nh *noteHandler) NoteList(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		return ErrNotFound
	}
	files := []string{
		"views/templates/base.html",
		"views/templates/pages/home.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		return ErrInternal
	}
	slog.Info("Executou o handler / ")
	return t.ExecuteTemplate(w, "base", nil)
}

func (nh *noteHandler) NoteView(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return apperror.WithStatus(errors.New("anotação é obrigatória"), http.StatusBadRequest)
	}
	if id == "0" {
		return apperror.WithStatus(errors.New("anotação 0 não foi encontrada"), http.StatusNotFound)
	}
	files := []string{
		"views/templates/base.html",
		"views/templates/pages/note-view.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		return ErrInternal
	}
	return t.ExecuteTemplate(w, "base", id)
}

func (nh *noteHandler) NoteNew(w http.ResponseWriter, r *http.Request) error {
	files := []string{
		"views/templates/base.html",
		"views/templates/pages/note-new.html",
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		return ErrInternal
	}
	return t.ExecuteTemplate(w, "base", nil)
}

func (nh *noteHandler) NoteCreate(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {

		w.Header().Set("Allow", http.MethodPost)

		//rejeitar a requisição
		return apperror.WithStatus(errors.New("operação não permitida"), http.StatusMethodNotAllowed)
	}
	fmt.Fprint(w, "Criando uma nova nota...")
	return nil
}
