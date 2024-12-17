package jwt

import (
	"errors"
	"fmt"
	"time"

	httpcommon "chi-mysql-boilerplate/internal/domain/http_common"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Payload interface{}
	jwt.RegisteredClaims
}

func GenerateToken(duration time.Duration, payload interface{}, isRefreshToken bool) (string, error) {
	// create a new JWT token
	claims := jwt.MapClaims{
		"exp":     time.Now().Add(duration).Unix(), // expires at
		"iat":     time.Now().Unix(),               // issued at
		"payload": payload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var key string
	if isRefreshToken {
		key = httpcommon.JwtConstants.RefreshSecretKey
	} else {
		key = httpcommon.JwtConstants.AccessSecretKey
	}
	// sign the token with the secret key
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(tokenString string, isRefreshToken bool) (*TokenClaims, error) {
	key := httpcommon.JwtConstants.AccessSecretKey
	if isRefreshToken {
		key = httpcommon.JwtConstants.RefreshSecretKey
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid %s token", getTokenType(isRefreshToken))
	}

	return claims, nil
}

func getTokenType(isRefreshToken bool) string {
	if isRefreshToken {
		return "refresh"
	}
	return "access"
}
