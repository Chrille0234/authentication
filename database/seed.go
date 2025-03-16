package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var query = `
    CREATE TABLE IF NOT EXISTS user (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      email varchar(128) UNIQUE NOT NULL,
      first_name varchar(64),
      last_name varchar(64),
      password_hashed varchar(64) NOT NULL
  );
  CREATE TABLE IF NOT EXISTS access_token (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER NOT NULL,
      token VARCHAR(64) UNIQUE NOT NULL,
      expires_at DATETIME NOT NULL,
      FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
  );

  CREATE TABLE IF NOT EXISTS refresh_token (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER NOT NULL,
      token VARCHAR(64) UNIQUE NOT NULL,
      expires_at DATETIME NOT NULL,
      FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
  );
`

func Seed(db *sqlx.DB) {
	db.MustExec(query)
}

func ConnectAndSeed(file string) *sqlx.DB {
	db := sqlx.MustConnect("sqlite3", file)
	Seed(db)
	return db
}
