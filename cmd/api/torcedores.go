package main

import (
	"encoding/json"
	"net/http"
	"github.com/sjakk/sjafoot/internal/data"

	"github.com/lib/pq"
)

func (app *application) registerTorcedorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nome  string `json:"nome"`
		Email string `json:"email"`
		Time  string `json:"time"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	torcedor := &data.Torcedor{
		Nome:      input.Nome,
		Email:     input.Email,
		TimeClube: input.Time,
	}

	v := data.NewValidator()
	v.Check(torcedor.Nome != "", "nome", "must be provided")
	v.Check(torcedor.Email != "", "email", "must be provided")
	v.Check(data.Matches(torcedor.Email, data.EmailRX), "email", "must be a valid email address")
	v.Check(torcedor.TimeClube != "", "time", "must be provided")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Torcedores.Insert(torcedor)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			v.AddError("email", "a fan with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := struct {
		ID       int64  `json:"id"`
		Nome     string `json:"nome"`
		Email    string `json:"email"`
		Time     string `json:"time"`
		Mensagem string `json:"mensagem"`
	}{
		ID:       torcedor.ID,
		Nome:     torcedor.Nome,
		Email:    torcedor.Email,
		Time:     torcedor.TimeClube,
		Mensagem: "Cadastro realizado com sucesso",
	}

	app.writeJSON(w, http.StatusCreated, response, nil)
}
