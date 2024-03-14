package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/robsondevgo/quicknotes/internal/repositories"
	"github.com/robsondevgo/quicknotes/utils"
)

type userHandler struct {
	repo repositories.UserRepository
}

func NewUserHandler(repo repositories.UserRepository) *userHandler {
	return &userHandler{repo: repo}
}

func (uh *userHandler) SignupForm(w http.ResponseWriter, r *http.Request) error {
	return render(w, http.StatusOK, "user-signup.html", nil)
}

func (uh *userHandler) Signup(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	data := newUserRequest(email, password)

	if strings.TrimSpace(data.Password) == "" {
		data.AddFieldError("password", "Senha é obrigatória")
	}

	if len(strings.TrimSpace(data.Password)) < 6 {
		data.AddFieldError("password", "Senha precisa ter no mínimo 6 caracteres")
	}

	if !isEmailValid(data.Email) || strings.TrimSpace(data.Password) == "" {
		data.AddFieldError("email", "Email é inválido")
	}

	if !data.Valid() {
		return render(w, http.StatusUnprocessableEntity, "user-signup.html", data)
	}

	//gerar o hash da senha
	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return err
	}
	hashToken := utils.GenerateTokenKey()
	user, token, err := uh.repo.Create(r.Context(), data.Email, hash, hashToken)
	if err == repositories.ErrDuplicateEmail {
		data.AddFieldError("email", "Email já está em uso")
		return render(w, http.StatusUnprocessableEntity, "user-signup.html", data)
	}

	if err != nil {
		return err
	}

	fmt.Println("Usuário criado:", user.Id)

	return render(w, http.StatusOK, "user-signup-success.html", token)
}

func (uh *userHandler) Confirm(w http.ResponseWriter, r *http.Request) error {
	token := r.PathValue("token")
	err := uh.repo.ConfirmUserByToken(r.Context(), token)
	msg := "Seu cadastro foi confirmado. Agora você já pode fazer o login no sistema"
	if err != nil {
		msg = "Esse cadastro já foi confirmado ou o token é inválido."
	}
	return render(w, http.StatusOK, "user-confirm.html", msg)
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
