package auth

import (
	"context"
	"fmt"
	"gcom/config"
	"gcom/types"
	"gcom/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

const UserKey string = "userID"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second + time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"userID": strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, userStore types.UserStore) http.HandlerFunc {
	return func (w http.ResponseWriter, r*http.Request) {
		tokenStr := getTokenFromRequest(r)
		token, err := validateToken(tokenStr)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return 
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)
		userID, _ := strconv.Atoi(str)
		u, err := userStore.GetUserById(userID)
		if err != nil {
			log.Println("failed to get user by id :%v", err)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if token != "" {
		return token
	}

	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return  -1
	}

	return userID
}