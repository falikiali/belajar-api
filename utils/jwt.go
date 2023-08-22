package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(id string) string {
	// Buat payload (klaim) untuk token Anda
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * 60).Unix(),
	}

	// Buat token JWT menggunakan claims dan kunci rahasia (secret key)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ganti "secretKey" dengan kunci rahasia yang lebih aman
	secretKey := []byte("secret-key-api-learn-golang")

	// Tanda tangani token menggunakan kunci rahasia
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		tokenString = ""
	}

	return tokenString
}

func CheckToken(tokenString string) (string, error) {
	secretKey := []byte("secret-key-api-learn-golang")

	// Parsing token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Cek error pada parsing token
	if err != nil {
		if err.Error() != "Token is expired" {
			return "", fmt.Errorf("Token is invalid")
		}

		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["id"].(string)
		return id, nil
	} else {
		return "", fmt.Errorf("Token is invalid")
	}
}
