package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TokenController struct {
	log *logrus.Entry
	ts  *TokenService
}

func NewTokenController(log *logrus.Logger, ts *TokenService) *TokenController {
	return &TokenController{
		log: log.WithField("controller", "token"),
		ts:  ts,
	}
}

func (tc *TokenController) HandleLogin(c *gin.Context) {
	tc.log.Info("handle login and generate token")
	form := UserClaims{}
	if err := c.ShouldBindJSON(&form); err != nil {
		tc.log.Error("failed ot bind JSON in handle login")
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
		tc.log.Error("missing param of user claim")
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing param",
			},
		)
		return
	}

	token, err := tc.ts.GenerateToken(form)
	if err != nil {
		tc.log.Error("failed to generate token")
		c.JSON(
			http.StatusInternalServerError,
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

func (tc *TokenController) HandleVerify(c *gin.Context) {
	tc.log.Info("handling verify token")
	Token := c.GetHeader("token")

	if Token == "" {
		tc.log.Error("missing mandatory input parameter in function HandleVerify")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing mandatory input parameter in function HandleVerify",
			},
		)
		return
	}

	claim, err := tc.ts.VerifyToken(Token)

	if err != nil {
		tc.log.Error("authentication error in handle verify")
		c.JSON(
			http.StatusInternalServerError,
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
			"claim": claim.Username,
		},
	)
}

func (tc *TokenController) HandleRefreshToken(c *gin.Context) {
	tc.log.Info("handle refresh token")
	Token := c.GetHeader("token")

	if Token == "" {
		tc.log.Error("missing mandatory input parameter in function HandleVerify")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing mandatory input parameter in function HandleVerify",
			},
		)
		return
	}

	claim, err := tc.ts.VerifyToken(Token)

	if err != nil {
		tc.log.Error("authentication error in handle refresh token")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	newToken, err := tc.ts.GenerateToken(*claim)

	if err != nil {
		tc.log.Error("failed to generate token")
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
			"token":    newToken,
			"username": claim.Username,
		},
	)
}
