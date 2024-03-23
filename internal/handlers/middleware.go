package handlers

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

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
