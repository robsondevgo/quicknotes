package main

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robsondevgo/quicknotes/internal/handlers"
	"github.com/robsondevgo/quicknotes/internal/mailer"
	"github.com/robsondevgo/quicknotes/internal/render"
	"github.com/robsondevgo/quicknotes/internal/repositories"
)

func LoadRoutes(sessionManager *scs.SessionManager, mail mailer.MailService, dbpool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	staticHandler := http.FileServer(http.Dir("views/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", staticHandler))

	noteRepo := repositories.NewNoteRepository(dbpool)
	userRepo := repositories.NewUserRepository(dbpool)

	render := render.NewRender(sessionManager)

	noteHandler := handlers.NewNoteHandler(render, sessionManager, noteRepo)
	userHandler := handlers.NewUserHandler(render, sessionManager, mail, userRepo)

	authMidd := handlers.NewAuthMiddleware(sessionManager)

	mux.HandleFunc("GET /", handlers.NewHomeHandler(render).HomeHandler)

	mux.Handle("GET /note", authMidd.RequireAuth(handlers.HandlerWithError(noteHandler.NoteList)))
	mux.Handle("GET /note/{id}", authMidd.RequireAuth(handlers.HandlerWithError(noteHandler.NoteView)))
	mux.Handle("GET /note/new", authMidd.RequireAuth(handlers.HandlerWithError(noteHandler.NoteNew)))
	mux.Handle("POST /note", authMidd.RequireAuth(handlers.HandlerWithError(noteHandler.NoteSave)))
	mux.Handle("DELETE /note/{id}", authMidd.RequireAuth(handlers.HandlerWithError(noteHandler.NoteDelete)))
	mux.Handle("GET /note/{id}/edit", authMidd.RequireAuth(handlers.HandlerWithError(noteHandler.NoteEdit)))

	mux.Handle("GET /user/signup", handlers.HandlerWithError(userHandler.SignupForm))
	mux.Handle("POST /user/signup", handlers.HandlerWithError(userHandler.Signup))

	mux.Handle("GET /user/signin", handlers.HandlerWithError(userHandler.SigninForm))
	mux.Handle("POST /user/signin", handlers.HandlerWithError(userHandler.Signin))

	mux.Handle("GET /user/signout", handlers.HandlerWithError(userHandler.Signout))

	mux.Handle("GET /user/forgetpassword", handlers.HandlerWithError(userHandler.ForgetPasswordForm))
	mux.Handle("POST /user/forgetpassword", handlers.HandlerWithError(userHandler.ForgetPassword))
	mux.Handle("POST /user/password", handlers.HandlerWithError(userHandler.ResetPassword))
	mux.Handle("GET /user/password/{token}", handlers.HandlerWithError(userHandler.ResetPasswordForm))

	mux.Handle("GET /me", authMidd.RequireAuth(handlers.HandlerWithError(userHandler.Me)))

	mux.Handle("GET /confirmation/{token}", handlers.HandlerWithError(userHandler.Confirm))

	return mux
}
