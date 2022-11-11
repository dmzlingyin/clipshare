package utils

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("set_your_secret")

type Claims struct {
	UserName string `json:"username"`
	Device   string `json:"device"`
	jwt.StandardClaims
}

func GenerateToken(username, device string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(72 * time.Hour)

	claims := Claims{
		username,
		device,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "clipshare",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
