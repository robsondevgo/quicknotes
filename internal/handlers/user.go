package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/robsondevgo/quicknotes/internal/render"
	"github.com/robsondevgo/quicknotes/internal/repositories"
	"github.com/robsondevgo/quicknotes/utils"
)

type userHandler struct {
	render  *render.RenderTemplate
	session *scs.SessionManager
	repo    repositories.UserRepository
}

func NewUserHandler(render *render.RenderTemplate, session *scs.SessionManager, repo repositories.UserRepository) *userHandler {
	return &userHandler{render: render, session: session, repo: repo}
}

func (uh *userHandler) Me(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Dados do usuário")
	return nil
}

func (uh *userHandler) SigninForm(w http.ResponseWriter, r *http.Request) error {
	return uh.render.RenderPage(w, r, http.StatusOK, "user-signin.html", nil)
}

func (uh *userHandler) Signin(w http.ResponseWriter, r *http.Request) error {
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

	if !isEmailValid(data.Email) {
		data.AddFieldError("email", "Email é inválido")
	}

	if !data.Valid() {
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signin.html", data)
	}

	//consultar o usuario pelo email
	user, err := uh.repo.FindByEmail(r.Context(), data.Email)
	if err != nil {
		data.AddFieldError("validation", "Credenciais inválidas")
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signin.html", data)
	}

	//verificar se o usuário está ativo
	if !user.Active.Bool {
		data.AddFieldError("validation", "Usuário não confirmou o cadastro")
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signin.html", data)
	}

	//validar se a senha está correta
	if !utils.ValidatePassword(data.Password, user.Password.String) {
		data.AddFieldError("validation", "Credenciais inválidas")
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signin.html", data)
	}

	//renew token
	err = uh.session.RenewToken(r.Context())
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	//armazena o id do usuário na sessão
	uh.session.Put(r.Context(), "userId", user.Id.Int.Int64())
	uh.session.Put(r.Context(), "userEmail", user.Email.String)

	http.Redirect(w, r, "/note", http.StatusSeeOther)
	return nil
}

func (uh *userHandler) SignupForm(w http.ResponseWriter, r *http.Request) error {
	return uh.render.RenderPage(w, r, http.StatusOK, "user-signup.html", nil)
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
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signup.html", data)
	}

	//gerar o hash da senha
	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return err
	}
	hashToken := utils.GenerateTokenKey()
	_, token, err := uh.repo.Create(r.Context(), data.Email, hash, hashToken)
	if err == repositories.ErrDuplicateEmail {
		data.AddFieldError("email", "Email já está em uso")
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signup.html", data)
	}

	if err != nil {
		return err
	}

	return uh.render.RenderPage(w, r, http.StatusOK, "user-signup-success.html", token)
}

func (uh *userHandler) Confirm(w http.ResponseWriter, r *http.Request) error {
	token := r.PathValue("token")
	err := uh.repo.ConfirmUserByToken(r.Context(), token)
	msg := "Seu cadastro foi confirmado. Agora você já pode fazer o login no sistema"
	if err != nil {
		msg = "Esse cadastro já foi confirmado ou o token é inválido."
	}
	return uh.render.RenderPage(w, r, http.StatusOK, "user-confirm.html", msg)
}

func (uh *userHandler) Signout(w http.ResponseWriter, r *http.Request) error {
	//renew token
	err := uh.session.RenewToken(r.Context())
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	uh.session.Remove(r.Context(), "userId")

	http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
	return nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
