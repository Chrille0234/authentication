package database

import (
	"fmt"
	"os"

	"github.com/chrille0234/auth/database/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

func CreateUser(user models.User, db *sqlx.DB) error {
	_, err := db.Exec("INSERT INTO user (first_name, last_name, password_hashed, email) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Password, user.Email)
	return err
}

func CreateToken(username string, expires int64) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      expires,
		})

	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetTokenInDB(token string, expires int64, user models.User, db *sqlx.DB) error {
	userFromDB, err := GetUserUsingEmail(user.Email, db)
	if err != nil {
		return err
	}

	query := `INSERT INTO access_token (token, expires_at, user_id) VALUES (?, ?, ?)`
	_, err = db.Exec(query, token, expires, userFromDB.ID)
	if err != nil {
		fmt.Printf("Error inserting token into database: %v", err)
		return err
	}

	return nil
}
