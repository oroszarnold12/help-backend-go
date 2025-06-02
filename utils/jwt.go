package utils

import (
	"fmt"
	"help/config"
	"help/constant"
	"help/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(user *model.User) (string, error) {
	expiration := time.Second * time.Duration(config.Env.JWTExpirationSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constant.UserUuidClaimKey: user.Uuid,
		"expiresAt":               time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Env.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("Cannot sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(config.Env.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("Cannot parse token: %w", err)
	}

	return token, nil
}
