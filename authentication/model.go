package authentication

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
