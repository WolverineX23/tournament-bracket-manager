package authentication

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type TokenService struct {
	log *logrus.Entry
}

func NewTokenService(log *logrus.Logger) *TokenService {
	return &TokenService{
		log: log.WithField("services", "Token"),
	}
}

var (
	secret     = []byte("16849841325189456f487")
	effectTime = 10 * time.Minute
)

func (ts *TokenService) GenerateToken(claims UserClaims) (string, error) {
	//设置token有效期，也可不设置有效期，采用redis的方式
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(effectTime).Unix()
	//生成token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		ts.log.Error("failed to generate token.")
		return "", err
	}
	return token, nil
}

func (ts *TokenService) VerifyToken(strToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		ts.log.Error("failed to parse with claims")
		return nil, errors.New("errors: authentication")
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		ts.log.Error("error in line 193: token.Claims.(*models.UserClaims)")
		return nil, errors.New("error: verify token")
	}
	if err := token.Claims.Valid(); err != nil {
		ts.log.Error("error in line 198: token.Claims.Valid()")
		return nil, errors.New("the claim is invalid in verify operation")
	}
	//fmt.Println("verify")
	return claims, nil
}
