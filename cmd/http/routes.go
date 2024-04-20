package main

import (
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robsondevgo/quicknotes/internal/handlers"
	"github.com/robsondevgo/quicknotes/internal/mailer"
	"github.com/robsondevgo/quicknotes/internal/render"
	"github.com/robsondevgo/quicknotes/internal/repositories"
	"github.com/robsondevgo/quicknotes/views"
)

func LoadRoutes(sessionManager *scs.SessionManager, mail mailer.MailService, dbpool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	static, err := fs.Sub(views.Files, "static")
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	staticHandler := http.FileServerFS(static)

	mux.Handle("GET /static/", http.StripPrefix("/static/", staticHandler))

	noteRepo := repositories.NewNoteRepository(dbpool)
	userRepo := repositories.NewUserRepository(dbpool)

	render := render.NewRender(sessionManager)

	noteHandler := handlers.NewNoteHandler(render, sessionManager, noteRepo)
	userHandler := handlers.NewUserHandler(render, sessionManager, mail, userRepo)

	authMidd := handlers.NewAuthMiddleware(sessionManager)
	errorMidd := handlers.NewErrorHandlerMiddleware(render)

	mux.HandleFunc("GET /", handlers.NewHomeHandler(render).HomeHandler)

	mux.Handle("GET /note", authMidd.RequireAuth(errorMidd.HandleError(noteHandler.NoteList)))
	mux.Handle("GET /note/{id}", authMidd.RequireAuth(errorMidd.HandleError(noteHandler.NoteView)))
	mux.Handle("GET /note/new", authMidd.RequireAuth(errorMidd.HandleError(noteHandler.NoteNew)))
	mux.Handle("POST /note", authMidd.RequireAuth(errorMidd.HandleError(noteHandler.NoteSave)))
	mux.Handle("DELETE /note/{id}", authMidd.RequireAuth(errorMidd.HandleError(noteHandler.NoteDelete)))
	mux.Handle("GET /note/{id}/edit", authMidd.RequireAuth(errorMidd.HandleError(noteHandler.NoteEdit)))

	mux.Handle("GET /user/signup", errorMidd.HandleError(userHandler.SignupForm))
	mux.Handle("POST /user/signup", errorMidd.HandleError(userHandler.Signup))

	mux.Handle("GET /user/signin", errorMidd.HandleError(userHandler.SigninForm))
	mux.Handle("POST /user/signin", errorMidd.HandleError(userHandler.Signin))

	mux.Handle("GET /user/signout", errorMidd.HandleError(userHandler.Signout))

	mux.Handle("GET /user/forgetpassword", errorMidd.HandleError(userHandler.ForgetPasswordForm))
	mux.Handle("POST /user/forgetpassword", errorMidd.HandleError(userHandler.ForgetPassword))
	mux.Handle("POST /user/password", errorMidd.HandleError(userHandler.ResetPassword))
	mux.Handle("GET /user/password/{token}", errorMidd.HandleError(userHandler.ResetPasswordForm))

	mux.Handle("GET /me", authMidd.RequireAuth(errorMidd.HandleError(userHandler.Me)))

	mux.Handle("GET /confirmation/{token}", errorMidd.HandleError(userHandler.Confirm))

	return mux
}
