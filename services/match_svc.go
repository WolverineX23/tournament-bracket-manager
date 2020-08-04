/*

 */

package services

import (
	"errors"
	"fmt"
	"os"

	"github.com/bitspawngg/tournament-bracket-manager/models"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type MatchService struct {
	log *logrus.Entry
	db  *models.DB
}

func getDB() (*models.DB, error) {
	db_type, exists := os.LookupEnv("DB_TYPE")
	if !exists {
		return nil, errors.New("missing DB_TYPE environment variable")
	}
	db_path, exists := os.LookupEnv("DB_PATH")
	if !exists {
		return nil, errors.New("missing DB_PATH environment variable")
	}
	db := models.NewDB(db_type, db_path)
	if err := db.Connect(); err != nil {
		return nil, err
	}
	return db, nil
}

func NewMatchService(log *logrus.Logger, db *models.DB) *MatchService {
	return &MatchService{
		log: log.WithField("services", "Match"),
		db:  db,
	}
}

func JudgeTeamNumber(n int) bool {
	if n <= 0 {
		return false
	}
	return n&(n-1) == 0
}

func GetMatchSchedule(teams []string, format string) ([]models.Match, string, error) {
	//根据传入队伍，初始化每场对战数据
	// implement proper check for number of teams in the next line
	if !JudgeTeamNumber(len(teams)) {
		return nil, "", errors.New("number of teams not a power of 2")
	}
	var matches []models.Match
	uuid4 := uuid.NewV4().String() //获取tournamentId的唯一值

	if format == "SINGLE" {
		lentemp := len(teams)
		round := 0
		for lentemp > 0 {
			lentemp /= 2
			round++
			if lentemp == len(teams)/2 {
				for i := 0; i < lentemp; i++ {
					matches = append(matches, models.Match{TournamentID: uuid4, Round: round, Table: i + 1, TeamOne: teams[2*i], TeamTwo: teams[2*i+1], Status: "Ready", Result: 0})
				}
			} else {
				for i := 0; i < lentemp; i++ {
					matches = append(matches, models.Match{TournamentID: uuid4, Round: round, Table: i + 1, TeamOne: "Unknown", TeamTwo: "Unknown", Status: "Pending", Result: 0})
				}
			}
		}

		db, err := getDB()
		if err != nil {
			return nil, "", err
		}
		if err := db.CreateMatches(matches); err != nil {
			return nil, "", err
		}
	} else if format == "CONSOLATION" {
		return nil, "", fmt.Errorf("Unsupported tournament format [%s]", format)
	} else {
		return nil, "", fmt.Errorf("Unsupported tournament format [%s]", format)
	}
	return matches, uuid4, nil
}

func SetMatchResult(tournamentId string, round, table, result int) error {
	if result != 1 && result != 2 {
		return errors.New("invalid result")
	}

	db, err := getDB()
	if err != nil {
		return err
	}

	pendingMatch, err := db.UpdateReadyMatch(tournamentId, round, table, result)

	if err != nil {
		return err
	}

	match := models.Match{}

	if err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, pendingMatch.TournamentID, pendingMatch.Round, pendingMatch.Table).First(&match).Error; err != nil {
		return nil //决赛后无比赛
	} else { //若当前比赛非决赛，则更新后续比赛数据
		if err := db.UpdatePendingMatch(pendingMatch); err != nil {
			return err
		}
		return nil
	}
}
