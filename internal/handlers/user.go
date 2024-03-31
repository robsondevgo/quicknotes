package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/robsondevgo/quicknotes/internal/mailer"
	"github.com/robsondevgo/quicknotes/internal/render"
	"github.com/robsondevgo/quicknotes/internal/repositories"
	"github.com/robsondevgo/quicknotes/utils"
)

type userHandler struct {
	render  *render.RenderTemplate
	session *scs.SessionManager
	mail    mailer.MailService
	repo    repositories.UserRepository
}

func NewUserHandler(render *render.RenderTemplate, session *scs.SessionManager, mail mailer.MailService, repo repositories.UserRepository) *userHandler {
	return &userHandler{render: render, session: session, mail: mail, repo: repo}
}

func (uh *userHandler) Me(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Dados do usuário")
	return nil
}

func (uh *userHandler) ForgetPasswordForm(w http.ResponseWriter, r *http.Request) error {
	return uh.render.RenderPage(w, r, http.StatusOK, "user-forget-password.html", nil)
}

func (uh *userHandler) ForgetPassword(w http.ResponseWriter, r *http.Request) error {
	//ler o email do formulário
	email := r.PostFormValue("email")

	//gerar um token
	hashToken := utils.GenerateTokenKey()

	//inserir um registro na tabela de tokens (user_confirmation_tokens)
	token, err := uh.repo.CreateResetPasswordToken(r.Context(), email, hashToken)

	if err != nil {
		data := UserRequest{}
		data.Email = email
		data.AddFieldError("email", "Email não possui cadastro válido no sistema")
		return uh.render.RenderPage(w, r, http.StatusOK, "user-forget-password.html", data)
	}

	//enviar um email com o link para atualizar a senha
	body, err := uh.render.RenderMailBody("forgetpassword.html", token)
	if err != nil {
		return err
	}
	err = uh.mail.Send(mailer.MailMessage{
		To:      []string{email},
		Subject: "Resetar senha",
		IsHtml:  true,
		Body:    body,
	})

	if err != nil {
		return err
	}

	message := "Foi enviado um email com um link para que você possa resetar a sua senha."

	return uh.render.RenderPage(w, r, http.StatusOK, "generic-success.html", message)
}

func (uh *userHandler) ResetPasswordForm(w http.ResponseWriter, r *http.Request) error {
	token := r.PathValue("token")

	userToken, err := uh.repo.GetUserConfirmationByToken(r.Context(), token)
	elapsedTime := time.Since(userToken.CreatedAt.Time).Hours()
	if err != nil || userToken.Confirmed.Bool || elapsedTime > 4 {
		msg := "Token inválido ou expirado. Solicite uma nova alteração."
		return uh.render.RenderPage(w, r, http.StatusOK, "generic-error.html", msg)
	}

	data := struct {
		Token  string
		Errors []string
	}{
		Token: token,
	}
	return uh.render.RenderPage(w, r, http.StatusOK, "user-reset-password.html", data)
}

func (uh *userHandler) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	//pegar os dados da senha
	password := r.PostFormValue("password")
	token := r.PostFormValue("token")

	//hash da senha
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		data := struct {
			Token  string
			Errors []string
		}{
			Token:  token,
			Errors: []string{"Não foi possível alterar a senha. Solicite uma nova alteração."},
		}
		return uh.render.RenderPage(w, r, http.StatusOK, "user-reset-password.html", data)
	}

	//atualizar a senha no banco
	email, err := uh.repo.UpdatePasswordByToken(r.Context(), hashedPassword, token)
	if err != nil {
		data := struct {
			Token  string
			Errors []string
		}{
			Token:  token,
			Errors: []string{"Não foi possível alterar a senha. Solicite uma nova alteração."},
		}
		return uh.render.RenderPage(w, r, http.StatusOK, "user-reset-password.html", data)
	}

	//enviar um email informando que a senha foi atualizada
	uh.mail.Send(mailer.MailMessage{
		To:      []string{email},
		Subject: "Sua senha foi atualizada",
		Body:    []byte("Sua senha foi atualizada e agora você já pode fazer o login novamente."),
	})

	uh.session.Put(r.Context(), "flash", "Sua senha foi atualizada. Agora você pode fazer o login.")

	http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
	return nil
}

func (uh *userHandler) SigninForm(w http.ResponseWriter, r *http.Request) error {
	data := UserRequest{}
	data.Flash = uh.session.PopString(r.Context(), "flash")
	return uh.render.RenderPage(w, r, http.StatusOK, "user-signin.html", data)
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
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signup.html", token)
	}

	if err != nil {
		return err
	}

	body, err := uh.render.RenderMailBody("confirmation.html", token)
	if err != nil {
		return err
	}

	//enviar email de confirmação do cadastro
	uh.mail.Send(mailer.MailMessage{
		To:      []string{data.Email},
		Subject: "Confirmação de Cadastro",
		IsHtml:  true,
		Body:    body,
	})

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
