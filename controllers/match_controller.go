package controllers

import (
	"net/http"

	"github.com/bitspawngg/tournament-bracket-manager/authentication"

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

func (mc *MatchController) HandlePing(c *gin.Context) {
	mc.log.Info("handling ping")

	form := authentication.Token{}

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

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, err := ts.VerifyToken(form.Token)

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

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, err := ts.VerifyToken(form.Token)

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
	var res[][] string;
	var temp[] string
	memo:=brackets[0].Round
	j:=-1;
	for i:=0;i< len(brackets);i++{
		if memo!=brackets[i].Round{
			j++
			memo=brackets[i].Round
			res=append(res,temp)
			temp=[]string{}
		}
		if brackets[i].TeamOne!="Unknown" {
			temp = append(temp, brackets[i].TeamOne)
		}else{
			temp = append(temp, "")
		}
		if brackets[i].TeamTwo!="Unknown" {
			temp = append(temp, brackets[i].TeamTwo)
		}else{
			temp = append(temp, "")
		}

		}
		if len(temp)!=0 {
			res=append(res,temp)
		}

	temp=[]string{}
	if brackets[len(brackets)-1].Result==1{
		temp=append(temp,brackets[len(brackets)-1].TeamOne)
		res=append(res,temp)
	}else {
		if brackets[len(brackets)-1].Result==2 {
			temp=append(temp,brackets[len(brackets)-1].TeamTwo)
			res=append(res,temp)
		}else {
			temp=append(temp,"")
			res=append(res,temp)
		}
	}

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

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, tErr := ts.VerifyToken(form.Token)

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
	err := mc.ms.SetMatchResultS(form.TournamentId, form.Round, form.Table, form.Result)
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

	brackets, tid, err := mc.ms.GetTour(form.TournamentId)
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
	var res[][] string;
	var temp[] string
	memo:=brackets[0].Round
	j:=-1;
	for i:=0;i< len(brackets);i++{
		if memo!=brackets[i].Round{
			j++
			memo=brackets[i].Round
			res=append(res,temp)
			temp=[]string{}
		}
		if brackets[i].TeamOne!="Unknown" {
			temp = append(temp, brackets[i].TeamOne)
		}else{
			temp = append(temp, "")
		}
		if brackets[i].TeamTwo!="Unknown" {
			temp = append(temp, brackets[i].TeamTwo)
		}else{
			temp = append(temp, "")
		}

	}
	if len(temp)!=0 {
		res=append(res,temp)
	}
	temp=[]string{}
	if brackets[len(brackets)-1].Result==1{
		temp=append(temp,brackets[len(brackets)-1].TeamOne)
		res=append(res,temp)
	}else {
		if brackets[len(brackets)-1].Result==2 {
			temp=append(temp,brackets[len(brackets)-1].TeamTwo)
			res=append(res,temp)
		}else {
			temp=append(temp,"")
			res=append(res,temp)
		}
	}


	c.JSON(
		http.StatusOK,
		gin.H{
			"data":          res,
			"tournament_id": tid,
			"username":      claim.Username,
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

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, tErr := ts.VerifyToken(form.Token)

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

	if form.TournamentId == "" || form.Round == 0 || form.Table == 0 {
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

	err := mc.ms.SetMatchResultC(form.TournamentId, form.Round, form.Table, form.Result)
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
	form := models.FormRefreshTable{}
	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Error("failed to bind JSON")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure 1",
				"error": err.Error(),
			},
		)
		return
	}

	log := authentication.ConfigureLogger()
	ts := authentication.NewTokenService(log)

	claim, tErr := ts.VerifyToken(form.Token)

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

	brackets, tid, err := mc.ms.GetTour(form.TournamentId)
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
	var res[][] string;
	var temp[] string
	memo:=brackets[0].Round
	j:=-1;
	for i:=0;i< len(brackets);i++{
		if memo!=brackets[i].Round{
			j++
			memo=brackets[i].Round
			res=append(res,temp)
			temp=[]string{}
		}
		if brackets[i].TeamOne!="Unknown" {
			temp = append(temp, brackets[i].TeamOne)
		}else{
			temp = append(temp, "winner")
		}
		if brackets[i].TeamTwo!="Unknown" {
			temp = append(temp, brackets[i].TeamTwo)
		}else{
			temp = append(temp, "winner")
		}

	}
	if len(temp)!=0 {
		res=append(res,temp)
	}
	temp=[]string{}
	if brackets[len(brackets)-1].Result==1{
		temp=append(temp,brackets[len(brackets)-1].TeamOne)
		res=append(res,temp)
	}else {
		if brackets[len(brackets)-1].Result==2 {
			temp=append(temp,brackets[len(brackets)-1].TeamTwo)
			res=append(res,temp)
		}else {
			temp=append(temp,"winner")
			res=append(res,temp)
		}
	}


	c.JSON(
		http.StatusOK,
		gin.H{
			"data":          res,
			"tournament_id": tid,
			"username":      claim.Username,
		},
	)
}