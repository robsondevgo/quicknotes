package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/robsondevgo/quicknotes/internal/handlers"
)

func main() {
	config := loadConfig()

	slog.SetDefault(newLogger(os.Stderr, config.GetLevelLog()))

	slog.Info(fmt.Sprintf("DB_PASSWORD: %s", config.DBPassword))

	slog.Info(fmt.Sprintf("Servidor rodando na porta %s\n", config.ServerPort))
	mux := http.NewServeMux()

	staticHandler := http.FileServer(http.Dir("views/static/"))

	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	noteHandler := handlers.NewNoteHandler()

	mux.HandleFunc("/", noteHandler.NoteList)
	mux.HandleFunc("/note/view", noteHandler.NoteView)
	mux.HandleFunc("/note/new", noteHandler.NoteNew)
	mux.HandleFunc("/note/create", noteHandler.NoteCreate)

	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), mux)
}
