package models

import "github.com/jmoiron/sqlx"

func DeleteExec(db *sqlx.DB, query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	return err
}

func IfExistExec(db *sqlx.DB, query string, args ...interface{}) bool {
	var res int
	_ = db.Get(&res, query, args...)
	return res != 0
}

func GetBy[T any](db *sqlx.DB, query string, args ...interface{}) (T, error) {
	var val T
	err := db.Get(&val, query, args...)
	return val, err
}

func SelectBy[T any](db *sqlx.DB, query string, args ...interface{}) (T, error) {
	var val T
	err := db.Select(&val, query, args...)
	return val, err
}
