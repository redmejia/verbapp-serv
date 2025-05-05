package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	BusinessName string `json:"business_name"`
	jwt.RegisteredClaims
}

func GenerateToken(jwtKey string, businessName string) (string, error) {

	exp := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		BusinessName: businessName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "", // app name
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(jwtKey)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func VerifyToken(token string, jwtKey string) (bool, *Claims, error) {

	secretKey := []byte(jwtKey)

	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return false, nil, err
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return false, nil, fmt.Errorf("invalid token")
	}

	return true, claims, nil
}
