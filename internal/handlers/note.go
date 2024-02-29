package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/robsondevgo/quicknotes/internal/apperror"
	"github.com/robsondevgo/quicknotes/internal/models"
	"github.com/robsondevgo/quicknotes/internal/repositories"
)

type noteHandler struct {
	repo repositories.NoteRepository
}

func NewNoteHandler(repo repositories.NoteRepository) *noteHandler {
	return &noteHandler{repo: repo}
}

func (nh *noteHandler) NoteList(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		return ErrNotFound
	}
	notes, err := nh.repo.List(r.Context())
	if err != nil {
		return err
	}
	return render(w, http.StatusOK, "home.html", newNoteResponseFromNoteList(notes))
}

func (nh *noteHandler) NoteView(w http.ResponseWriter, r *http.Request) error {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		return apperror.WithStatus(errors.New("anotação é obrigatória"), http.StatusBadRequest)
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}
	note, err := nh.repo.GetById(r.Context(), id)
	if err != nil {
		return err
	}
	return render(w, http.StatusOK, "note-view.html", newNoteResponseFromNote(note))
}

func (nh *noteHandler) NoteNew(w http.ResponseWriter, r *http.Request) error {
	return render(w, http.StatusOK, "note-new.html", newNoteRequest(nil))
}

func (nh *noteHandler) NoteSave(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {

		w.Header().Set("Allow", http.MethodPost)

		//rejeitar a requisição
		return apperror.WithStatus(errors.New("operação não permitida"), http.StatusMethodNotAllowed)
	}
	err := r.ParseForm()
	if err != nil {
		return err
	}
	idParam := r.PostForm.Get("id")
	id, _ := strconv.Atoi(idParam)
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	color := r.PostForm.Get("color")

	var note *models.Note
	if id > 0 {
		note, err = nh.repo.Update(r.Context(), id, title, content, color)
	} else {
		note, err = nh.repo.Create(r.Context(), title, content, color)
	}
	if err != nil {
		return err
	}
	http.Redirect(w, r, fmt.Sprintf("/note/view?id=%d", note.Id.Int), http.StatusSeeOther)
	return nil
}

func (nh *noteHandler) NoteDelete(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodPost)
		//rejeitar a requisição
		return apperror.WithStatus(errors.New("operação não permitida"), http.StatusMethodNotAllowed)
	}

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		return apperror.WithStatus(errors.New("anotação é obrigatória"), http.StatusBadRequest)
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}
	err = nh.repo.Delete(r.Context(), id)
	if err != nil {
		return err
	}
	return nil
}

func (nh *noteHandler) NoteEdit(w http.ResponseWriter, r *http.Request) error {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		return apperror.WithStatus(errors.New("anotação é obrigatória"), http.StatusBadRequest)
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}
	note, err := nh.repo.GetById(r.Context(), id)
	if err != nil {
		return err
	}
	return render(w, http.StatusOK, "note-edit.html", newNoteRequest(note))
}
