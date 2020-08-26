/*

 */

package models

import (
	"errors"

	_ "github.com/satori/go.uuid"
)

type Match struct {
	ID 			 uint   `json:"id" gorm:"primary_key"`
	TournamentID string `json:"tournamentId"`
	Round        int    `json:"round"`
	Table        int    `json:"table"`
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
func (db DB) UpdateReadyMatchS(tournamentId string, round, table, result int) (Match, error) {
	updateMatch := Match{}
	match := Match{}
	pendingMatch := Match{}
	zeroMatch := Match{}

	if err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, tournamentId, round, table).First(&match).Error; err != nil {
		return zeroMatch, err
	}

	if match.Result != 0 {
		return zeroMatch, errors.New("This match has finished")
	}

	if match.Status == "Pending" {
		return zeroMatch, errors.New("This match is Pending")
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

func (db DB) UpdatePendingMatchS(pendingMatch Match) error {
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

func (db DB) UpdateReadyMatchC(tournamentId string, round, table, result int) (Match, Match, error) {
	zeroMatch := Match{}
	updateMatch := Match{}
	winner := Match{}
	loser := Match{}
	match := Match{}
	oneRound := []Match{}

	if err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, tournamentId, round, table).First(&match).Error; err != nil {
		return zeroMatch, zeroMatch, err
	}

	/*var tablesOfOneRound int
	if err := db.DB.Where(`"tournament_id" = ? AND "round" = ?`, tournamentId, 1).Count(&tablesOfOneRound).Error; err != nil {
		return zeroMatch, zeroMatch, err
	}
	Wrong: out=0
	*/

	if err := db.DB.Where(`"tournament_id" = ? AND "round" = ?`, tournamentId, 1).Find(&oneRound).Error; err != nil {
		return zeroMatch, zeroMatch, err
	}
	tablesOfOneRound := len(oneRound)
	//print(tablesOfOneRound)

	if match.Result != 0 {
		return zeroMatch, zeroMatch, errors.New("This match has finished")
	}

	if match.Status == "Pending" {
		return zeroMatch, zeroMatch, errors.New("This match is pending")
	}

	updateMatch.Result = result
	updateMatch.Status = "Finished"

	if err := db.DB.Model(&match).Updates(updateMatch).Error; err != nil { //更新已完成的原ready比赛数据
		return zeroMatch, zeroMatch, err
	}

	//loser
	loser.TournamentID = tournamentId
	loser.Round = round + 1
	loser.Table = (table + 1) / 2
	if table%2 != 0 { //teamOne
		if result == 1 {
			loser.TeamOne = match.TeamTwo
		} else {
			loser.TeamOne = match.TeamOne
		}
	} else { //TeamTwo
		if result == 1 {
			loser.TeamTwo = match.TeamTwo
		} else {
			loser.TeamTwo = match.TeamOne
		}
	}

	//winner
	winner.TournamentID = tournamentId
	winner.Round = round + 1
	winner.Table = (table+1)/2 + tablesOfOneRound/2
	if table%2 != 0 { //teamOne
		if result == 1 {
			winner.TeamOne = match.TeamOne
		} else {
			winner.TeamOne = match.TeamTwo
		}
	} else { //TeamTwo
		if result == 1 {
			winner.TeamTwo = match.TeamOne
		} else {
			winner.TeamTwo = match.TeamTwo
		}
	}

	return winner, loser, nil
}

func (db DB) UpdatePendingMatchC(winner, loser Match) error {
	winnerMatch := Match{}
	loserMatch := Match{}

	//update winner
	if err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, winner.TournamentID, winner.Round, winner.Table).First(&winnerMatch).Error; err != nil {
		return err
	}
	if winnerMatch.TeamOne != "Unknown" || winnerMatch.TeamTwo != "Unknown" {
		winner.Status = "Ready"
	}
	if err := db.DB.Model(&winnerMatch).Updates(winner).Error; err != nil {
		return err
	}

	//update loser
	if err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, loser.TournamentID, loser.Round, loser.Table).First(&loserMatch).Error; err != nil {
		return err
	}
	if loserMatch.TeamOne != "Unknown" || loserMatch.TeamTwo != "Unknown" {
		loser.Status = "Ready"
	}
	if err := db.DB.Model(&loserMatch).Updates(loser).Error; err != nil {
		return err
	}

	return nil
}
