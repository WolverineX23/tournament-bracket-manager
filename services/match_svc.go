/*

 */

package services

import (
	"errors"
	"fmt"

	"github.com/bitspawngg/tournament-bracket-manager/models"
	"github.com/sirupsen/logrus"
)

type MatchService struct {
	log *logrus.Entry
	db  *models.DB
}

func NewMatchService(log *logrus.Logger, db *models.DB) *MatchService {
	return &MatchService{
		log: log.WithField("services", "Match"),
		db:  db,
	}
}

func GetMatchSchedule(teams []string, format string) ([]models.Match, error) {
	// implement proper check for number of teams in the next line
	if len(teams) == 3 {
		return nil, errors.New("number of teams not a power of 2")
	}
	var matches []models.Match
	if format == "SINGLE" {
		// To be implemented, remove code below
		matches = []models.Match{
			models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 1, Table: 1, TeamOne: "A", TeamTwo: "B", Status: "Ready", Result: 0},
			models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 1, Table: 2, TeamOne: "C", TeamTwo: "D", Status: "Ready", Result: 0},
			models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 2, Table: 1, TeamOne: "A", TeamTwo: "C", Status: "Pending", Result: 0},
		}
	} else if format == "CONSOLATION" {
		return nil, fmt.Errorf("Unsupported tournament format [%s]", format)
	} else {
		return nil, fmt.Errorf("Unsupported tournament format [%s]", format)
	}
	return matches, nil
}
