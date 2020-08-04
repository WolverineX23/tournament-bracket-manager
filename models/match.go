/*

 */

package models

import (
	"errors"

	_ "github.com/satori/go.uuid"
)

type Match struct {
	TournamentID string `json:"tournamentId" gorm:"primarykey"`
	Round        int    `json:"round" gorm:"primarykey"`
	Table        int    `json:"table" gorm:"primarykey"`
	TeamOne      string `json:"teamOne"`
	TeamTwo      string `json:"teamTwo"`
	Status       string `json:"status"`
	Result       int    `json:"result"` // 1 if Team One wins, 2 if Team Two wins, -1 if no winner
}

func (db DB) CreateMatches(matches []Match) error {
	return db.DB.Create(matches).Error
}

func (db DB) GetMatch(tournamentId string, round, table int) (*Match, error) {
	match := Match{}
	err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, tournamentId, round, table).First(&match).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (db DB) GetMatchesByTournament(tournamentId string) ([]Match, error) {
	matches := make([]Match, 0)
	err := db.DB.Order("round").Where("tournament_id = ?", tournamentId).Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (db DB) GetMatchesByStatus(status string) ([]Match, error) {
	matches := make([]Match, 0)
	err := db.DB.Where("status = ?", status).Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (db DB) DeleteMatch(tournamentId string, round, table int) error {
	match := Match{}
	err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, tournamentId, round, table).First(&match).Error
	if err != nil {
		return err
	}
	err = db.DB.Delete(&match).Error
	if err != nil {
		return err
	}
	return nil
}

//当输入一个对战结果，更新该组match的result and status,并返回关联到的PendingMatch
func (db DB) UpdateReadyMatch(tournamentId string, round, table, result int) (Match, error) {
	updateMatch := Match{}
	match := Match{}
	pendingMatch := Match{}
	zeroMatch := Match{}

	if err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, tournamentId, round, table).First(&match).Error; err != nil {
		return zeroMatch, err
	}

	if match.Result != 0 {
		return zeroMatch, errors.New("This tournament has finished")
	}
	
	if match.Status == "Pending" {
		return zeroMatch, errors.New("This tournament is Pending")
	}

	updateMatch.Result = result
	updateMatch.Status = "Finshed"

	if err := db.DB.Model(&match).Updates(updateMatch).Error; err != nil {
		return zeroMatch, err
	}

	pendingMatch.TournamentID = tournamentId
	pendingMatch.Round = round + 1
	if table%2 == 0 { //teamTwo
		pendingMatch.Table = table / 2
		if result == 1 {
			pendingMatch.TeamTwo = match.TeamOne
		} else {
			pendingMatch.TeamTwo = match.TeamTwo
		}
	} else { //TeamOne
		pendingMatch.Table = table/2 + 1
		if result == 1 {
			pendingMatch.TeamOne = match.TeamOne
		} else {
			pendingMatch.TeamOne = match.TeamTwo
		}
	}
	return pendingMatch, nil
}

func (db DB) UpdatePendingMatch(pendingMatch Match) error {
	match := Match{}

	if err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, pendingMatch.TournamentID, pendingMatch.Round, pendingMatch.Table).First(&match).Error; err != nil {
		return err
	}

	if match.TeamOne != "Unknown" || match.TeamTwo != "Unknown" {
		pendingMatch.Status = "Ready"
	}

	if err := db.DB.Model(&match).Updates(pendingMatch).Error; err != nil {
		return err
	}

	return nil
}
