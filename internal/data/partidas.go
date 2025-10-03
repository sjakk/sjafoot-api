package data

type Partida struct {
	TimeCasa string `json:"time_casa"`
	TimeFora string `json:"time_fora"`
	Placar   string `json:"placar"`
}

type PartidasResponse struct {
	Rodada   int       `json:"rodada"`
	Partidas []Partida `json:"partidas"`
}

type MatchesResponse struct {
	Matches []struct {
		HomeTeam struct {
			Name string `json:"name"`
		} `json:"homeTeam"`
		AwayTeam struct {
			Name string `json:"name"`
		} `json:"awayTeam"`
		Score struct {
			FullTime struct {
				HomeTeam int `json:"home"`
				AwayTeam int `json:"away"`
			} `json:"fullTime"`
		} `json:"score"`
		Matchday int `json:"matchday"`
	} `json:"matches"`
}
