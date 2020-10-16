package controllers

import (
	"net/http"

	"github.com/bitspawngg/tournament-bracket-manager/models"

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

func (mc *MatchController) Authenticate(c *gin.Context, token string) *models.UserClaims {

	claim, err := mc.ms.VerifyToken(token)
	if err != nil {
		mc.log.Error("failed to verify token")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return nil
	}

	return claim
}

func (mc *MatchController) HandlePing(c *gin.Context) {
	mc.log.Info("handling ping")

	form := models.Token{}

	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Error("failed to bind JSON in handle verify")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	if form.Token == "" {
		mc.log.Error("missing mandatory input parameter in function HandleVerify")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing mandatory input parameter in function HandleVerify",
			},
		)
		return
	}

	claim := mc.Authenticate(c, form.Token)

	if claim == nil {
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

func (mc *MatchController) HandleLogin(c *gin.Context) {
	mc.log.Info("handle login and generate token")
	form := models.UserClaims{}
	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Error("failed ot bind JSON in handle login")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	if form.ID == "" || form.Name == "" || form.Username == "" || form.Password == "" {
		mc.log.Error("missing param of user claim")
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing param",
			},
		)
		return
	}

	token, err := mc.ms.GenerateToken(form)
	if err != nil {
		mc.log.Error("failed to generate token")
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
			"msg":   "success",
			"token": token,
		},
	)
}

func (mc *MatchController) HandleVerify(c *gin.Context) {
	mc.log.Info("handling verify token")
	form := models.Token{}

	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Error("failed to bind JSON in hadnle verify")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	if form.Token == "" {
		mc.log.Error("missing mandatory input parameter in function HandleVerify")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing mandatory input parameter in function HandleVerify",
			},
		)
		return
	}

	claim := mc.Authenticate(c, form.Token)

	if claim == nil {
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"msg":   "success",
			"claim": claim.Username,
		},
	)
}

func (mc *MatchController) HandleRefreshToken(c *gin.Context) {
	mc.log.Info("handle refresh token")
	form := models.Token{}

	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Error("failed to bind JSON in hadnle verify")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	if form.Token == "" {
		mc.log.Error("missing mandatory input parameter in function HandleVerify")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing mandatory input parameter in function HandleVerify",
			},
		)
		return
	}

	claim := mc.Authenticate(c, form.Token)

	if claim == nil {
		return
	}

	newToken, err := mc.ms.GenerateToken(*claim)

	if err != nil {
		mc.log.Error("failed to generate token")
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
			"msg":   "success",
			"token": newToken,
		},
	)
}

func (mc *MatchController) HandleGetMatchSchedule(c *gin.Context) {
	mc.log.Info("handling match schedule")
	form := models.FormGetMatchSchedule{}
	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Error("failed to bind JSON in handle get match schedule")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	claim := mc.Authenticate(c, form.Token)

	if claim == nil {
		return
	}

	if form.Teams == nil {
		mc.log.Error("missing mandatory input parameter")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing mandatory input parameter",
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

func (mc *MatchController) HandleSetMatchResultS(c *gin.Context) {
	mc.log.Info("handing set match result of single")
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

	claim := mc.Authenticate(c, form.Token)

	if claim == nil {
		return
	}

	if form.TournamentId == "" || form.Round == 0 || form.Table == 0 {
		mc.log.Error("missing param")
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing param",
			},
		)
		return
	}

	err := mc.ms.SetMatchResultS(form.TournamentId, form.Round, form.Table, form.Result)
	if err != nil {
		mc.log.Error("failed to set match result")
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

func (mc *MatchController) HandleSetMatchResultC(c *gin.Context) {
	mc.log.Info("handing set match result of consolation")
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

	claim := mc.Authenticate(c, form.Token)

	if claim == nil {
		return
	}

	if form.TournamentId == "" || form.Round == 0 || form.Table == 0 {
		mc.log.Error("missing param")
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing param",
			},
		)
		return
	}

	err := mc.ms.SetMatchResultC(form.TournamentId, form.Round, form.Table, form.Result)
	if err != nil {
		mc.log.Error("failed to set match result")
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
