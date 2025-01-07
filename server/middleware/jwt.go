package middleware

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var (
	SecretKey = []byte("KsS2X1CgFT4bi3BRRIxLk5jjiUBj8wxE")
)

type JWTHandler struct {
}

type UserClaims struct {
	ID   int64
	Role string
	jwt.StandardClaims
}

func (j *JWTHandler) SetToken(id int64, role string) (string, error) {
	claims := UserClaims{
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Issuer:    "github.com/crazyfrankie/onlinejudge",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(SecretKey)

	return token, err
}
