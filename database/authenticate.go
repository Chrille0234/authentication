package database

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUserPassword(email, password string, db *sqlx.DB) (bool, error) {
	user, err := GetUserUsingEmail(email, db)

	if err != nil {
		return false, err
	}

	isValidPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil

	if !isValidPassword {
		return false, nil
	}

	return true, nil
}

func AuthenticateToken(tokenString string) (bool, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return false, err
	}

	// Verify the token isn't expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return false, fmt.Errorf("token expired")
			}
			return true, nil
		}
	}

	return false, fmt.Errorf("invalid token claims")
}
