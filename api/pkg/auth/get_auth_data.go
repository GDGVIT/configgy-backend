package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/GDGVIT/configgy-backend/api/pkg/services/authsvc"
	"github.com/golang-jwt/jwt/v4"
)

func GetAuthorizationHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	return authHeader, nil
}

func GetAuthDataFromToken(signedToken string) (*authsvc.AuthData, error) {
	var res *authsvc.AuthData
	jwtKey := []byte(os.Getenv("TOKEN_ACCESS_SECRET"))
	token, err := jwt.ParseWithClaims(signedToken, &authsvc.AuthData{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return res, err
	}
	res, ok := token.Claims.(*authsvc.AuthData)
	if !ok {
		err = errors.New("couldn't parse claims")
		return res, err
	}
	return res, err
}

func GetRefreshTokenDataFromToken(signedToken string) (*authsvc.AuthData, error) {
	var res *authsvc.AuthData
	jwtKey := []byte(os.Getenv("TOKEN_REFRESH_SECRET"))
	token, err := jwt.ParseWithClaims(signedToken, &authsvc.AuthData{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return res, err
	}
	res, ok := token.Claims.(*authsvc.AuthData)
	if !ok {
		err = errors.New("couldn't parse claims")
		return res, err
	}
	return res, err
}

func ValidateToken(signedToken string) error {
	accessSecretKey := []byte(os.Getenv("TOKEN_ACCESS_SECRET"))
	token, err := jwt.ParseWithClaims(signedToken, &authsvc.AuthData{}, func(t *jwt.Token) (interface{}, error) {
		return accessSecretKey, nil
	})
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*authsvc.AuthData)
	if !ok {
		err = errors.New("couldn't parse claims")
		return err
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return err
	}

	// extra verification
	_, err = jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_ACCESS_SECRET")), nil
	})
	if err != nil {
		return err
	}

	return err
}
