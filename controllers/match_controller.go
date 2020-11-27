package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/bitspawngg/tournament-bracket-manager/authentication"

	"github.com/bitspawngg/tournament-bracket-manager/models"
	"github.com/bitspawngg/tournament-bracket-manager/services"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

type MatchController struct {
	log          *logrus.Entry
	ms           *services.MatchService
	socketServer *socketio.Server
}

func NewMatchController(log *logrus.Logger, ms *services.MatchService, socketServer *socketio.Server) *MatchController {
	return &MatchController{
		log:          log.WithField("controller", "match"),
		ms:           ms,
		socketServer: socketServer,
	}
}

func (mc *MatchController) HandlePing(c *gin.Context) {
	mc.log.Info("handling ping")

	Token := c.GetHeader("token")

	if Token == "" {
		mc.log.Error("anthentication error:token is null in handling 'ping' require")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "token is null",
			},
		)
		return
	}

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, err := ts.VerifyToken(Token)

	if err != nil {
		mc.log.Error("Authentication error in handle ping")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"msg":   "pong",
			"claim": claim.Username,
		},
	)

}

func (mc *MatchController) HandleGetMatchSchedule(c *gin.Context) {
	mc.log.Info("handling match schedule")
	Token := c.GetHeader("token")

	if Token == "" {
		mc.log.Error("anthentication error:token is null in handling 'matchschedule' require")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "token is null",
			},
		)
		return
	}

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, err := ts.VerifyToken(Token)

	if err != nil {
		mc.log.Error("Authentication error in handle get match schedule")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	form := models.FormGetMatchSchedule{}
	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Error("failed to bind JSON in handle get match schedule")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":      "failure",
				"error":    err.Error(),
				"username": claim.Username,
			},
		)
		return
	}

	if form.Teams == nil {
		mc.log.Error("missing mandatory input parameter")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":      "failure",
				"error":    "missing mandatory input parameter",
				"username": claim.Username,
			},
		)
		return
	}

	brackets, tid, err := mc.ms.GetMatchSchedule(form.Teams, form.Format)
	if err != nil {
		mc.log.Error("failed to get match schedule")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"msg":      "failure",
				"error":    err.Error(),
				"username": claim.Username,
			},
		)
		return
	}

	res := mc.ms.GetResult(brackets)

	c.JSON(
		http.StatusOK,
		gin.H{
			"data":          res,
			"tournament_id": tid,
			"username":      claim.Username,
		},
	)
}

func (mc *MatchController) HandleSetMatchResultS(c *gin.Context) {
	mc.log.Info("handing set match result of single")
	Token := c.GetHeader("token")
	TournamentID := c.Param("tournamentID")

	if Token == "" || TournamentID == "" {
		mc.log.Error("anthentication error:token or tournamentID is null in handling 'results' require")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "token or tournamentID is null",
			},
		)
		return
	}

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, tErr := ts.VerifyToken(Token)

	if tErr != nil {
		mc.log.Error("Authentication error in handle set result of single match")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": tErr.Error(),
			},
		)
		return
	}

	form := models.FormSetMatchResult{}
	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Error("failed to bind JSON")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":      "failure",
				"error":    err.Error(),
				"username": claim.Username,
			},
		)
		return
	}

	err := mc.ms.SetMatchResultS(TournamentID, form.Round, form.Table, form.Result)
	if err != nil {
		mc.log.Error("failed to set match result")
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg":      "failure",
				"error":    err.Error(),
				"username": claim.Username,
			},
		)
		return
	}

	round := form.Round + 1
	table := (form.Round + 1) / 2
	teamName, err := mc.ms.GetWinTeam(TournamentID, form.Round, form.Table)
	if err != nil {
		mc.log.Error("failed to get win team")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"msg":      "failure",
				"error":    err.Error(),
				"username": claim.Username,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"round":         round,
			"tournament_id": TournamentID,
			"table":         table,
			"team_name":     teamName,
			"username":      claim.Username,
		},
	)

	inform := models.FormWinner{
		TournamentId: TournamentID,
		Round:        round,
		Table:        table,
		TeamName:     teamName,
	}

	data, err := json.Marshal(inform)
	if err != nil {
		mc.log.Error("failed to marshal json for inform")
	}
	// broadcast
	mc.socketServer.BroadcastToRoom("", "tournament", "update", data)

}

func (mc *MatchController) HandleSetMatchResultC(c *gin.Context) {
	mc.log.Info("handing set match result of consolation")
	Token := c.GetHeader("token")
	TournamentID := c.Param("tournamentID")

	if Token == "" || TournamentID == "" {
		mc.log.Error("anthentication error:token or tournamentID is null in handling 'resultc' require")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "token or tournamentID is null",
			},
		)
		return
	}

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, tErr := ts.VerifyToken(Token)

	if tErr != nil {
		mc.log.Error("Authentication error in handle set result of consolation match")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": tErr.Error(),
			},
		)
		return
	}

	form := models.FormSetMatchResult{}

	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Error("failed to bind JSON")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	if TournamentID == "" || form.Round == 0 || form.Table == 0 {
		mc.log.Error("missing param")
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":      "failure",
				"error":    "missing param",
				"username": claim.Username,
			},
		)
		return
	}

	err := mc.ms.SetMatchResultC(TournamentID, form.Round, form.Table, form.Result)
	if err != nil {
		mc.log.Error("failed to set match result")
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg":      "failure",
				"error":    err.Error(),
				"username": claim.Username,
			},
		)
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"msg":      "success",
			"username": claim.Username,
		},
	)
}

func (mc *MatchController) HandleRefreshTable(c *gin.Context) {
	mc.log.Info("handing set match result of single")
	Token := c.GetHeader("token")
	TournamentID := c.Param("tournamentID")

	if Token == "" || TournamentID == "" {
		mc.log.Error("anthentication error:token or tournamentID is null in handling 'matchschedule' require")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "token or tournamentID is null",
			},
		)
		return
	}

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, tErr := ts.VerifyToken(Token)

	if tErr != nil {
		mc.log.Error("Authentication error in handle set result of single match")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure 2",
				"error": tErr.Error(),
			},
		)
		return
	}

	brackets, err := mc.ms.GetTour(TournamentID)
	if err != nil {
		mc.log.Error("failed to get match schedule")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"msg":      "failure",
				"error":    err.Error(),
				"username": claim.Username,
			},
		)
		return
	}

	res := mc.ms.GetResult(brackets)

	c.JSON(
		http.StatusOK,
		gin.H{
			"data":          res,
			"tournament_id": TournamentID,
			"username":      claim.Username,
		},
	)
}

func (mc *MatchController) HandleGetAlltournamentID(c *gin.Context) {
	mc.log.Info("handing get all tournamentId")
	Token := c.GetHeader("token")

	if Token == "" {
		mc.log.Error("anthentication error:token is null in handling 'alltournament' require")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "token is null",
			},
		)
		return
	}

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, tErr := ts.VerifyToken(Token)

	if tErr != nil {
		mc.log.Error("Authentication error in handling getAllTounamentID.")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure 2",
				"error": tErr.Error(),
			},
		)
		return
	}

	tournamentId, err := mc.ms.GetAllTourID()
	if err != nil {
		mc.log.Error("failed to get tournamentId")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error":    err,
				"username": claim.Username,
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"data":     tournamentId,
			"username": claim.Username,
		},
	)

}

func (mc *MatchController) HandleGetRate(c *gin.Context) {
	mc.log.Info("handing get rate")
	Token := c.GetHeader("token")
	TournamentID := c.Param("tournamentID")

	if Token == "" || TournamentID == "" {
		mc.log.Error("anthentication error:token or tournamentID is null in handling 'matchschedule' require")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "token or tournamentID is null",
			},
		)
		return
	}

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, tErr := ts.VerifyToken(Token)

	if tErr != nil {
		mc.log.Error("Authentication error in handle getRate")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure 2",
				"error": tErr.Error(),
			},
		)
		return
	}

	rate, err := mc.ms.GetRateOfWinning(TournamentID)
	if err != nil {
		mc.log.Error("failed to get rate of winning")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error":    err,
				"username": claim.Username,
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"data":     rate,
			"username": claim.Username,
		},
	)

}
