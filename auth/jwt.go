package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	TokenType   string `json:"token_type"`
	jwt.StandardClaims
}

func GenerateJWT(userID, phoneNumber, tokenType string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	if tokenType == "refresh" {
		expirationTime = time.Now().Add(7 * 24 * time.Hour) // Refresh token lasts longer
	}

	claims := &Claims{
		UserID:      userID,
		PhoneNumber: phoneNumber,
		TokenType:   tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
