package models

type FormGetMatchSchedule struct {
	Teams  []string `json:"teams"`
	Format string   `json:"format"`
}

type FormSetMatchResult struct {
	Round  int `json:"round"`
	Table  int `json:"table"`
	Result int `json:"result"`
}

type FormWinner struct {
	TournamentId string `json:"tournament_id"`
	Round        int    `json:"round"`
	Table        int    `json:"table"`
	TeamName     string `json:"team_name"`
}
