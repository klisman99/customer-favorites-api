package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtTokenService struct {
	secretKey  string
	expireTime time.Duration
}

func NewTokenService(secretKey string, expireTime time.Duration) *jwtTokenService {
	return &jwtTokenService{secretKey: secretKey, expireTime: expireTime}
}

func (s *jwtTokenService) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.expireTime).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *jwtTokenService) Validate(token string) (string, error) {
	token, err := parseTokenPrefix(token)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid token")
	}

	return userID, nil
}

func parseTokenPrefix(token string) (string, error) {
	if len(token) < 7 || token[:7] != "Bearer " {
		return "", errors.New("invalid token format")
	}
	return token[7:], nil
}
