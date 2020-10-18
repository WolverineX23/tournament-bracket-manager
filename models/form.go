package models

type FormGetMatchSchedule struct {
	Token  string   `json:"token"`
	Teams  []string `json:"teams"`
	Format string   `json:"format"`
}

type FormSetMatchResult struct {
	Token        string `json:"token"`
	TournamentId string `json:"tournament_id"`
	Round        int    `json:"round"`
	Table        int    `json:"table"`
	Result       int    `json:"result"`
}
