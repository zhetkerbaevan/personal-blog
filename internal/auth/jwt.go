package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/zhetkerbaevan/personal-blog/internal/config"
	"github.com/zhetkerbaevan/personal-blog/internal/models"
	"github.com/zhetkerbaevan/personal-blog/internal/utils"
)

type contextKey string
const UserKey contextKey = "userId"

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

func WithJWTAuth(handlerFunc http.HandlerFunc, store models.UserStoreInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get token from request
		tokenString := getTokenFromRequest(r)

		//validate JWT
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("FAILED TO VALIDATE TOKEN %v", err)
			permissionDenied(w)
			return
		}
		if !token.Valid {
			log.Println("INVALID TOKEN")
			permissionDenied(w)
			return
		}

		//extract claims from token
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userId"].(string) //get userId 

		userId, _ := strconv.Atoi(str) //convert it to int

	
		u, err := store.GetUserByID(userId)
		if err != nil {
			log.Printf("FAILED TO GET USER BY ID %v", err)
			permissionDenied(w)
			return
		}

		//add userId to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.Id) //create child context, which has key-value ("userId" : id)
		r = r.WithContext(ctx) //return new request with context (creates new copy of request and returns with updated context)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" {
		return tokenAuth
	}
	return ""
}

func validateJWT(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) { //validate token
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(config.Envs.JWTSecret), nil
		}
		return nil, fmt.Errorf("UNEXPECTED SIGNING METHOD %v", t.Header["alg"])
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("PERMISSON DENIED"))
}

func GetUserIdFromContext(ctx context.Context) int {
	userId, ok := ctx.Value(UserKey).(int) //return value of key
	if ok {
		return userId
	}
	return -1
}