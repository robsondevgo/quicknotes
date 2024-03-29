package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robsondevgo/quicknotes/internal/mailer"
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

	//testar o envio de email
	mailService := mailer.NewSMTPMailService(mailer.SMTPConfig{
		Host:     "localhost",
		Port:     1025,
		Username: "",
		Password: "",
		From:     "quicknotes@quick.com",
	})

	mailService.Send(mailer.MailMessage{
		To:      []string{"user1@quick.com"},
		Subject: "Email de teste",
		IsHtml:  true,
		Body:    []byte("<h1>Esta é uma mensagem de teste</h1>"),
	})

	sessionManager := scs.New()
	sessionManager.Lifetime = time.Hour
	sessionManager.Store = pgxstore.New(dbpool)
	//Limpa as sessões expiradas da tabela de sessions a cada 30 minutos
	pgxstore.NewWithCleanupInterval(dbpool, 30*time.Second)

	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))

	mux := LoadRoutes(sessionManager, dbpool)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), sessionManager.LoadAndSave(csrfMiddleware(mux))); err != nil {
		panic(err)
	}
}
