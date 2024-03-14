package handlers

import (
	"fmt"

	"github.com/robsondevgo/quicknotes/internal/models"
	"github.com/robsondevgo/quicknotes/internal/validations"
)

type NoteResponse struct {
	Id      int
	Title   string
	Content string
	Color   string
}

type NoteRequest struct {
	Id      int
	Title   string
	Content string
	Color   string
	Colors  []string
	validations.FormValidator
}

type UserRequest struct {
	Email    string
	Password string
	validations.FormValidator
}

func newUserRequest(email, password string) (req UserRequest) {
	req.Email = email
	req.Password = password
	return
}

func newNoteRequest(note *models.Note) (req NoteRequest) {
	for i := 1; i <= 9; i++ {
		req.Colors = append(req.Colors, fmt.Sprintf("color%d", i))
	}
	if note != nil {
		req.Id = int(note.Id.Int.Int64())
		req.Title = note.Title.String
		req.Color = note.Color.String
		req.Content = note.Content.String
	} else {
		req.Color = "color3"
	}
	return
}

func newNoteResponseFromNote(note *models.Note) (res NoteResponse) {
	res.Id = int(note.Id.Int.Int64())
	res.Title = note.Title.String
	res.Content = note.Content.String
	res.Color = note.Color.String
	return
}

func newNoteResponseFromNoteList(notes []models.Note) (res []NoteResponse) {
	for _, note := range notes {
		res = append(res, newNoteResponseFromNote(&note))
	}
	return
}
