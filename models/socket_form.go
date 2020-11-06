package models

type FormRefresh struct {
	Status   string `json:"status"`
	Round    int    `json:"round"`
	Table    int    `json:"table"`
	TeamName string `json:"team_name"`
}

type FormLeave struct {
	TournamentId string
}
