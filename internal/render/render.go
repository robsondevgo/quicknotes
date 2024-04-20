package render

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/csrf"
	"github.com/robsondevgo/quicknotes/views"
)

type RenderTemplate struct {
	session *scs.SessionManager
}

func NewRender(session *scs.SessionManager) *RenderTemplate {
	return &RenderTemplate{session: session}
}

func getTemplatePageFiles(t *template.Template, page string, useFS bool) (*template.Template, error) {
	if useFS {
		return t.ParseFS(views.Files, "templates/base.html", "templates/pages/"+page)
	}
	files := []string{
		"views/templates/base.html",
	}
	files = append(files, "views/templates/pages/"+page)
	return t.ParseFiles(files...)
}

func getTemplateMailFiles(mailTmpl string, useFS bool) (*template.Template, error) {
	if useFS {
		return template.ParseFS(views.Files, "templates/mails/"+mailTmpl)
	}
	return template.ParseFiles("views/templates/mails/" + mailTmpl)
}

func (rt *RenderTemplate) RenderPage(w http.ResponseWriter, r *http.Request, status int, page string, data any) error {
	t := template.New("").Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"csrfToken": func() string {
			return csrf.Token(r)
		},
		"isAuthenticated": func() bool {
			return rt.session.Exists(r.Context(), "userId")
		},
		"userEmail": func() string {
			return rt.session.GetString(r.Context(), "userEmail")
		},
	})
	useFS := !strings.Contains(r.Host, "localhost")
	t, err := getTemplatePageFiles(t, page, useFS)
	if err != nil {
		return err
	}
	buff := &bytes.Buffer{}
	err = t.ExecuteTemplate(buff, "base", data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	buff.WriteTo(w)
	return nil
}

func (rt *RenderTemplate) RenderMailBody(r *http.Request, mailTempl string, data map[string]string) ([]byte, error) {
	useFS := !strings.Contains(r.Host, "localhost")

	//TODO: verificar se está usando HTTPS
	data["hostAddr"] = "http://" + r.Host

	t, err := getTemplateMailFiles(mailTempl, useFS)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	w := &bytes.Buffer{}
	err = t.Execute(w, data)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return w.Bytes(), nil
}
