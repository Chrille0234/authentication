package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `db:"id"`
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Password  string `db:"password_hashed"`
}

func (u *User) HashAndSalt() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}
