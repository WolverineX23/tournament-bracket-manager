/*

 */

package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/bitspawngg/tournament-bracket-manager/models"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type MatchService struct {
	log *logrus.Entry
	db  *models.DB
}

var (
	secret     = []byte("16849841325189456f487")
	effectTime = 1 * time.Minute
)

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

func (ms *MatchService) GenerateToken(claims models.UserClaims) (string, error) {
	//设置token有效期，也可不设置有效期，采用redis的方式
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(effectTime).Unix()
	//生成token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		ms.log.Error("failed to generate token.")
		return "", err
	}
	return token, nil
}

func (ms *MatchService) VerifyToken(strToken string) (*models.UserClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		ms.log.Error("failed to parse with claims")
		return nil, err
	}
	claims, ok := token.Claims.(*models.UserClaims)
	if !ok {
		ms.log.Error("error in line 193: token.Claims.(*models.UserClaims)")
		return nil, errors.New("error: verify token")
	}
	if err := token.Claims.Valid(); err != nil {
		ms.log.Error("error in line 198: token.Claims.Valid()")
		return nil, errors.New("the claim is invalid in verify operation")
	}
	fmt.Println("verify")
	return claims, nil
}
