package handlers

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/robsondevgo/quicknotes/internal/apperror"
	"github.com/robsondevgo/quicknotes/internal/repositories"
)

var ErrNotFound = apperror.WithStatus(errors.New("página não encontrada"), http.StatusNotFound)
var ErrInternal = apperror.WithStatus(errors.New("aconteceu um erro ao executar essa operação"), http.StatusInternalServerError)

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func (f HandlerWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		var statusErr apperror.StatusError
		var repoErr repositories.RepositoryError
		if errors.As(err, &statusErr) {
			if statusErr.StatusCode() == http.StatusNotFound {
				files := []string{
					"views/templates/base.html",
					"views/templates/pages/404.html",
				}
				t, err := template.ParseFiles(files...)
				if err != nil {
					http.Error(w, err.Error(), statusErr.StatusCode())
				}
				t.ExecuteTemplate(w, "base", statusErr.Error())
				return
			}
			http.Error(w, err.Error(), statusErr.StatusCode())
			return
		}
		if errors.As(err, &repoErr) {
			slog.Error(err.Error())
			http.Error(w, "aconteceu um erro ao executar essa operação", http.StatusInternalServerError)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
