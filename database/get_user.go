package database

import (
	"github.com/chrille0234/auth/database/models"
	"github.com/jmoiron/sqlx"
)

func GetUserUsingEmail(email string, db *sqlx.DB) (models.User, error) {
	user := models.User{}
	err := db.Get(&user, "SELECT * FROM user WHERE email = $1", email)

	return user, err
}

func GetUserFromToken(token string, db *sqlx.DB) (models.User, error) {
	user := models.User{}
	err := db.Get(&user, `SELECT user.* FROM user
    JOIN access_token ON user.id = access_token.user_id
    WHERE access_token.token = $1`, token)

	return user, err
}
