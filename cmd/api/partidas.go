package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/sjakk/sjafoot/internal/data"
	"strconv"
)

func (app *application) listPartidasHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	url := fmt.Sprintf("http://api.football-data.org/v4/competitions/%d/matches", id)

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

	var externalResponse data.MatchesResponse
	if err := json.NewDecoder(resp.Body).Decode(&externalResponse); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	equipe := r.URL.Query().Get("equipe")
	rodadaStr := r.URL.Query().Get("rodada")

	var filteredMatches []data.Partida

	for _, match := range externalResponse.Matches {
		if rodadaStr != "" {
			rodada, _ := strconv.Atoi(rodadaStr)
			if match.Matchday != rodada {
				continue
			}
		}

		if equipe != "" {
			if match.HomeTeam.Name != equipe && match.AwayTeam.Name != equipe {
				continue
			}
		}

		// If it passes all filters, add it to the list
		filteredMatches = append(filteredMatches, data.Partida{
			TimeCasa:  match.HomeTeam.Name,
			TimeFora:  match.AwayTeam.Name,
			Placar:    fmt.Sprintf("%d-%d", match.Score.FullTime.HomeTeam, match.Score.FullTime.AwayTeam),
		})
	}

	rodadaFinal, _ := strconv.Atoi(rodadaStr)
	response := data.PartidasResponse{
		Rodada:   rodadaFinal,
		Partidas: filteredMatches,
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
