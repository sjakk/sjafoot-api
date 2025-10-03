package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) broadcastHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Tipo     string `json:"tipo"`
		Time     string `json:"time"`
		Placar   string `json:"placar,omitempty"`
		Mensagem string `json:"mensagem"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Tipo != "inicio" && input.Tipo != "fim" {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, "invalid 'tipo' field")
		return
	}

	torcedores, err := app.models.Torcedores.GetAllForTeam(input.Time)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, torcedor := range torcedores {
		logMessage := fmt.Sprintf("BROADCAST TO: %s (%s) | MSG: %s", torcedor.Nome, torcedor.Email, input.Mensagem)
		app.logger.Println(logMessage)
	}

	response := map[string]interface{}{
		"status":          "broadcast initiated",
		"team":            input.Time,
		"event_type":      input.Tipo,
		"notified_fans": len(torcedores),
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}
