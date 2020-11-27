/*

 */

package services

import (
	"errors"
	"fmt"
	"math"

	"github.com/bitspawngg/tournament-bracket-manager/models"
	uuid "github.com/satori/go.uuid"
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

func JudgeTeamNumber(n int) bool {
	if n <= 0 {
		return false
	}
	return n&(n-1) == 0
}

func GetMatchCount(n int) int {
	count := 0
	i := 0
	temp := n

	for temp != 1 {
		temp >>= 1
		i++
	}

	count = n / 2 * i
	return count
}

func (ms *MatchService) GetMatchSchedule(teams []string, format string) ([]models.Match, string, error) {
	//根据传入队伍，初始化每场对战数据
	// implement proper check for number of teams in the next line
	if !JudgeTeamNumber(len(teams)) {
		ms.log.Error("number og teams not a power of 2")
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

		if err := ms.db.CreateMatches(matches); err != nil {
			ms.log.Error("failed to create matches")
			return nil, "", err
		}
		ms.log.Info("Create single tournament successfully")
	} else if format == "CONSOLATION" {
		count := GetMatchCount(len(teams)) //获取安慰赛总比赛场次count
		round := 0
		sub := len(teams) / 2

		for count != 0 {
			count -= sub
			round++
			if round == 1 {
				for i := 0; i < sub; i++ {
					matches = append(matches, models.Match{TournamentID: uuid4, Round: round, Table: i + 1, TeamOne: teams[2*i], TeamTwo: teams[2*i+1], Status: "Ready", Result: 0})
				}
			} else {
				for i := 0; i < sub; i++ {
					matches = append(matches, models.Match{TournamentID: uuid4, Round: round, Table: i + 1, TeamOne: "Unknown", TeamTwo: "Unknown", Status: "Pending", Result: 0})
				}
			}
		}

		if err := ms.db.CreateMatches(matches); err != nil {
			ms.log.Error("failed to create matches")
			return nil, "", err
		}
		ms.log.Info("Create consolation tournament successfully")
	} else {
		ms.log.Error("Unsupported tournament format")
		return nil, "", fmt.Errorf("Unsupported tournament format [%s]", format)
	}

	return matches, uuid4, nil
}

func (ms *MatchService) SetMatchResultS(tournamentId string, round, table, result int) error {
	if result != 1 && result != 2 {
		ms.log.Error("input an invalid result")
		return errors.New("invalid result")
	}

	pendingMatch, err := ms.db.UpdateReadyMatchS(tournamentId, round, table, result)

	if err != nil {
		ms.log.Error("failed to update ready match of single")
		return err
	}

	match := models.Match{}

	if err := ms.db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, pendingMatch.TournamentID, pendingMatch.Round, pendingMatch.Table).First(&match).Error; err != nil {
		return nil //决赛后无比赛
	} else { //若当前比赛非决赛，则更新后续比赛数据
		if err := ms.db.UpdatePendingMatchS(pendingMatch); err != nil {
			ms.log.Error("failed to update penging match of single")
			return err
		}
		return nil
	}
}

func (ms *MatchService) SetMatchResultC(tournamentId string, round, table, result int) error {
	if result != 1 && result != 2 {
		ms.log.Error("input an invalid result")
		return errors.New("invalid result")
	}

	winner, loser, err := ms.db.UpdateReadyMatchC(tournamentId, round, table, result)

	if err != nil {
		ms.log.Error("failed to update ready match of consolation")
		return err
	}

	match := models.Match{}
	if err := ms.db.DB.Where(`"tournament_id" = ? AND "round" = ?`, tournamentId, round+1).First(&match).Error; err != nil {
		return nil
	} else {
		if err := ms.db.UpdatePendingMatchC(winner, loser); err != nil {
			ms.log.Error("failed to update penging match of consolation")
			return err
		}
		return nil
	}
}

func (ms *MatchService) GetTour(id string) ([]models.Match, error) {
	matches, err := ms.db.GetMatchesByTournament(id)
	return matches, err
}

func (ms *MatchService) GetResult(matches []models.Match) [][]string {
	var res [][]string
	var temp []string
	memo := matches[0].Round
	j := -1
	for i := 0; i < len(matches); i++ {
		if memo != matches[i].Round {
			j++
			memo = matches[i].Round
			res = append(res, temp)
			temp = []string{}
		}
		if matches[i].TeamOne != "Unknown" {
			temp = append(temp, matches[i].TeamOne)
		} else {
			temp = append(temp, "")
		}
		if matches[i].TeamTwo != "Unknown" {
			temp = append(temp, matches[i].TeamTwo)
		} else {
			temp = append(temp, "")
		}

	}
	if len(temp) != 0 {
		res = append(res, temp)
	}
	temp = []string{}
	if matches[len(matches)-1].Result == 1 {
		temp = append(temp, matches[len(matches)-1].TeamOne)
		res = append(res, temp)
	} else {
		if matches[len(matches)-1].Result == 2 {
			temp = append(temp, matches[len(matches)-1].TeamTwo)
			res = append(res, temp)
		} else {
			temp = append(temp, "")
			res = append(res, temp)
		}
	}

	return res
}

func (ms *MatchService) GetWinTeam(tournamentId string, round, table int) (string, error) {
	winner, err := ms.db.GetWinTeamByStatus(tournamentId, round, table)
	return winner, err
}

func (ms *MatchService) GetAllTourID() ([]string, error) {
	var tournamentId []string
	matches := make([]models.Match, 0)
	err := ms.db.DB.Find(&matches).Error
	if err != nil {
		return nil, err
	}
	flag := false
	for i := 0; i < len(matches); i++ {
		flag = false
		for j := 0; j < len(tournamentId); j++ {
			if matches[i].TournamentID == tournamentId[j] {
				flag = true
				break
			}
		}
		if flag == false {
			tournamentId = append(tournamentId, matches[i].TournamentID)
		}
	}
	return tournamentId, nil
}

func (ms *MatchService) GetChampion(tournamentId string) (string, error) {
	matches, err := ms.db.GetMatchesByTournament(tournamentId)
	if err != nil {
		return "", err
	}
	thisMatch, err := ms.db.GetMatch(tournamentId, int(math.Log2(float64(len(matches)))+1), 1)
	if err != nil {
		return "", err
	}
	if thisMatch == nil {
		return "", nil
	}
	if thisMatch.Result == 0 {
		return "", nil
	}
	if thisMatch.Result == 1 {
		return thisMatch.TeamOne, nil
	}
	return thisMatch.TeamTwo, nil

}

func (ms *MatchService) GetTeams(tournamentId string) ([]string, error) {
	matches, err := ms.db.GetMatchesByTournament(tournamentId)
	if err != nil {
		return nil, err
	}
	flag1 := false
	flag2 := false
	var teams []string
	for i := 0; i < len(matches); i++ {
		flag1 = false
		flag2 = false
		for j := 0; j < len(teams); j++ {
			if matches[i].TeamOne == teams[j] {
				flag1 = true
			}
			if matches[i].TeamTwo == teams[j] {
				flag2 = true
			}
		}
		if flag1 == false {
			teams = append(teams, matches[i].TeamOne)
		}
		if flag2 == false {
			teams = append(teams, matches[i].TeamTwo)
		}
	}
	return teams, nil
}
func (ms *MatchService) GetRateOfWinning(tournamentId string) ([]float64, error) {
	tournamentIds, err := ms.GetAllTourID()
	if err != nil {
		return nil, err
	}
	Teams, err := ms.GetTeams(tournamentId)
	if err != nil {
		return nil, err
	}
	var rates []float64
	var cntOfWinning int
	cntOfWinning = 0
	for j := 0; j < len(Teams); j++ {
		// 获取每个队的胜率
		cntOfWinning = 0
		for i := 0; i < len(tournamentIds); i++ {
			// 遍历每场比赛
			winner, err := ms.GetChampion(tournamentIds[i]) // 获取该比赛的胜利队伍
			if err != nil {
				return nil, err
			}
			if winner == Teams[j] {
				cntOfWinning++
			}
		}
		rate := float64(cntOfWinning) / float64(len(tournamentIds))
		rate *= 100
		rate = math.Trunc(rate*1e2+0.5) * 1e-2
		rates = append(rates, rate)
		rate = 0
	}
	return rates, nil
}
