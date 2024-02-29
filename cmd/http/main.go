package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robsondevgo/quicknotes/internal/handlers"
	"github.com/robsondevgo/quicknotes/internal/repositories"
)

func main() {
	config := loadConfig()

	slog.SetDefault(newLogger(os.Stderr, config.GetLevelLog()))

	dbpool, err := pgxpool.New(context.Background(), config.DBConnURL)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Conexão com o banco aconteceu com sucesso")

	defer dbpool.Close()

	slog.Info(fmt.Sprintf("Servidor rodando na porta %s\n", config.ServerPort))
	mux := http.NewServeMux()

	staticHandler := http.FileServer(http.Dir("views/static/"))

	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	noteRepo := repositories.NewNoteRepository(dbpool)

	noteHandler := handlers.NewNoteHandler(noteRepo)

	mux.Handle("/", handlers.HandlerWithError(noteHandler.NoteList))
	mux.Handle("/note/view", handlers.HandlerWithError(noteHandler.NoteView))
	mux.Handle("/note/new", handlers.HandlerWithError(noteHandler.NoteNew))
	mux.Handle("/note/save", handlers.HandlerWithError(noteHandler.NoteSave))
	mux.Handle("/note/delete", handlers.HandlerWithError(noteHandler.NoteDelete))
	mux.Handle("/note/edit", handlers.HandlerWithError(noteHandler.NoteEdit))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), mux); err != nil {
		panic(err)
	}
}
