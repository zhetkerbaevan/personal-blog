package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/zhetkerbaevan/personal-blog/internal/config"
)

func CreateJWT(secret []byte, userId int) (string, error) {
	//set expiration time 
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ //using HMAC-SHA256 for signing method, and arguments for claims
		"userId" : strconv.Itoa(userId), 
		"expiredAt" : time.Now().Add(expiration).Unix(), // seconds starting from 1 january 1970 to (now + expiration)
	})

	tokenString, err := token.SignedString(secret) //signing token if true returns jwt
	if err != nil {
		return "", err
	}
	return tokenString, nil
}