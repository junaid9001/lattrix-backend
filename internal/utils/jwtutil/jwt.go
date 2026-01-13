package jwtutil

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/junaid9001/lattrix-backend/internal/config"
)

type AccessClaims struct {
	Role        string `json:"role"`
	WorkspaceID string `json:"workspace_id"`
	TokenType   string `json:"token_type"`

	jwt.RegisteredClaims
}

type RefreshClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func CreateAccessToken(userID int, workSpaceID string, role string) (string, error) {
	secret := config.AppConfig.JWT_KEY

	claims := AccessClaims{
		Role:        role,
		WorkspaceID: workSpaceID,
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			Issuer:    "lattrix",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(240 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ValidateAccessToken(tokenStr string) (*AccessClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token not found")
	}

	secret := config.AppConfig.JWT_KEY

	claims := &AccessClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != "access" {
		return nil, errors.New("not an access token")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenStr string) (*RefreshClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token not found")
	}

	secret := config.AppConfig.JWT_KEY

	claims := &RefreshClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("not a refresh token")
	}

	return claims, nil
}
func CreateRefreshToken(userID int) (string, error) {
	secret := config.AppConfig.JWT_KEY

	claims := RefreshClaims{
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			Issuer:    "lattrix",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}
