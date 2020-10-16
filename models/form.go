package models

import (
	"github.com/dgrijalva/jwt-go"
)

//用户信息类，作为生成token的参数
type UserClaims struct {
	jwt.StandardClaims
	ID       string `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type FormGetMatchSchedule struct {
	Token  string   `json:"token"`
	Teams  []string `json:"teams"`
	Format string   `json:"format"`
}

type FormSetMatchResult struct {
	Token        string `json:"token"`
	TournamentId string `json:"tournament_id"`
	Round        int    `json:"round"`
	Table        int    `json:"table"`
	Result       int    `json:"result"`
}

type Token struct {
	Token string `json:"token"`
}
