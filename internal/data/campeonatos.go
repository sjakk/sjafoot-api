package data

type Campeonato struct {
	ID        int    `json:"id"`
	Nome      string `json:"nome"`
	Temporada string `json:"temporada"`
}

type CompetitionsResponse struct {
	Competitions []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		CurrentSeason struct {
			StartDate string `json:"startDate"`
		} `json:"currentSeason"`
	} `json:"competitions"`
}
