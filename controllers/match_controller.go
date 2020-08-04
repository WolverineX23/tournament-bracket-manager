package controllers

import (
	"net/http"

	"github.com/bitspawngg/tournament-bracket-manager/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MatchController struct {
	log *logrus.Entry
	ms  *services.MatchService
}

func NewMatchController(log *logrus.Logger, ms *services.MatchService) *MatchController {
	return &MatchController{
		log: log.WithField("controller", "match"),
		ms:  ms,
	}
}

func (mc *MatchController) HandlePing(c *gin.Context) {
	mc.log.Info("handling ping")
	c.JSON(
		http.StatusOK,
		gin.H{
			"msg": "pong",
		},
	)
}

type FormGetMatchSchedule struct {
	Teams  []string `json:"teams"`
	Format string   `json:"format"`
}

type FormSetMatchResult struct {
	TournamentId string `json:"tournament_id"`
	Round        int    `json:"round"`
	Table        int    `json:"table"`
	Result       int    `json:"result"`
}

func (mc *MatchController) HandleGetMatchSchedule(c *gin.Context) {
	form := FormGetMatchSchedule{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	if form.Teams == nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing mandatory input parameter",
			},
		)
		return
	}

	brackets, tid, err := services.GetMatchSchedule(form.Teams, form.Format)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
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
			"msg":           "success",
			"data":          brackets,
			"tournament_id": tid,
		},
	)
}

func (mc *MatchController) HandleSetMatchResult(c *gin.Context) {
	form := FormSetMatchResult{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	if form.TournamentId == "" || form.Round == 0 || form.Table == 0 {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing param",
			},
		)
		return
	}
	
	err := services.SetMatchResult(form.TournamentId, form.Round, form.Table, form.Result)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"msg": "success",
		},
	)
}
