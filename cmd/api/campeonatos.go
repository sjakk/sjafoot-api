package main

import (
	"encoding/json"
	"net/http"
	"github.com/sjakk/sjafoot/internal/data"
)

func (app *application) listCampeonatosHandler(w http.ResponseWriter, r *http.Request) {
	url := "http://api.football-data.org/v4/competitions/"

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	req.Header.Set("X-Auth-Token", "6311a66f5f8746fd8860a5de6173f49f")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var externalResponse data.CompetitionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&externalResponse); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Transform the data to the desired format
	campeonatos := make([]data.Campeonato, 0, len(externalResponse.Competitions))
	for _, comp := range externalResponse.Competitions {
		campeonatos = append(campeonatos, data.Campeonato{
			ID:        comp.ID,
			Nome:      comp.Name,
			Temporada: comp.CurrentSeason.StartDate[:4],
		})
	}

	err = app.writeJSON(w, http.StatusOK, campeonatos, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
