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

type FormRefreshTable struct {
	Token        string `json:"token"`
	TournamentId string `json:"tournament_id"`
}

type FormWinner struct {
	TournamentId string `json:"tournament_id"`
	Round        int    `json:"round"`
	Table        int    `json:"table"`
	TeamName     string `json:"team_name"`
}

type FormToken struct {
	Token string `json:"token"`
}
type FormGetRate struct {
	Token string `json:"token"`
	Team  string `json:"team"`
}
