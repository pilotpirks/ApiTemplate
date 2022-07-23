package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id       string    `json:"id" db:"id"`
	Name     string    `validate:"required" json:"name" db:"name"`
	Password string    `validate:"required" json:"password" db:"password"`
	JwtKey   string    `json:"jwt_key" db:"jwt_key"`
	Created  time.Time `json:"created" db:"created"`
}

func AddUser(val *User, db *sqlx.DB) error {
	_, err := db.Exec("INSERT INTO users (id, name, password, jwt_key, created) VALUES ($1, $2, $3, $4, $5)",
		val.Id, val.Name, val.Password, val.JwtKey, time.Now())
	return err
}

func UpdateUser(val *User, db *sqlx.DB) error {
	_, err := db.NamedExec(`UPDATE users SET
		name = :name,
		password = :password,
		jwt_key = :jwt_key
			WHERE id = :id`,
		map[string]interface{}{
			"name":     val.Name,
			"password": val.Password,
			"jwt_key":  val.JwtKey,
			"id":       val.Id,
		})

	return err
}
