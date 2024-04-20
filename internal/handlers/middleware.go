package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/robsondevgo/quicknotes/internal/apperror"
	"github.com/robsondevgo/quicknotes/internal/render"
	"github.com/robsondevgo/quicknotes/internal/repositories"
)

var ErrNotFound = apperror.WithStatus(errors.New("página não encontrada"), http.StatusNotFound)
var ErrInternal = apperror.WithStatus(errors.New("aconteceu um erro ao executar essa operação"), http.StatusInternalServerError)

type authMiddleware struct {
	session *scs.SessionManager
}

func NewAuthMiddleware(session *scs.SessionManager) *authMiddleware {
	return &authMiddleware{session: session}
}

func (ah *authMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := ah.session.GetInt64(r.Context(), "userId")
		if userId == 0 {
			slog.Warn("usuário não está logado")
			http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type errorHandlerMiddleware struct {
	render *render.RenderTemplate
}

func NewErrorHandlerMiddleware(render *render.RenderTemplate) *errorHandlerMiddleware {
	return &errorHandlerMiddleware{render: render}
}

func (em *errorHandlerMiddleware) HandleError(next func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			var statusErr apperror.StatusError
			var repoErr repositories.RepositoryError
			if errors.As(err, &statusErr) {
				if statusErr.StatusCode() == http.StatusNotFound {
					//renderizar uma página de erro 404
					em.render.RenderPage(w, r, http.StatusNotFound, "404.html", nil)
					return
				}
			}

			//repositories errors
			if errors.As(err, &repoErr) {
				slog.Error(err.Error())
				em.render.RenderPage(w, r, http.StatusInternalServerError, "generic-error.html", "aconteceu um erro ao executar essa operação")
				return
			}

			//all other generic errors
			slog.Error(err.Error())
			em.render.RenderPage(w, r, http.StatusInternalServerError, "generic-error.html", err.Error())
		}
	})
}
