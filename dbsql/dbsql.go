package dbsql

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// ------------------------------------------------

type iDatabase interface {
	DbConnect(param string) (*sqlx.DB, error)
	DbClose(db *sqlx.DB)
}

// ------------------------------------------------

type database struct {
	dbType string
}

func (d *database) DbConnect(param string) (*sqlx.DB, error) {
	db, err := sqlx.Open(d.dbType, param)
	return db, err
}

func (d *database) DbClose(db *sqlx.DB) {
	db.Close()
}

// ------------------------------------------------

// factory
func GetDatabase(t string) (iDatabase, error) {
	if t == "sqlite3" {
		return newSqLite(), nil
	}
	if t == "mysql" {
		return newMySql(), nil
	}
	if t == "postgres" {
		return newPostgres(), nil
	}

	return nil, errors.New("Wrong database type")
}

// ------------------------------------------------

type sqlite_ struct {
	database
}

func newSqLite() iDatabase {
	return &sqlite_{
		database: database{
			dbType: "sqlite3",
		},
	}
}

// ------------------------------------------------

type mysql_ struct {
	database
}

func newMySql() iDatabase {
	return &mysql_{
		database: database{
			dbType: "mysql",
		},
	}
}

// ------------------------------------------------

type postgres_ struct {
	database
}

func newPostgres() iDatabase {
	return &postgres_{
		database: database{
			dbType: "postgres",
		},
	}
}

// ------------------------------------------------

/*
	sqlite: "./adm.db"
	mysql: "test:test@(localhost:3306)/test"
	postgres: "user=foo dbname=bar sslmode=disable"

	tst, err := dbsql.GetDatabase("sqlite3")
	if err != nil {
		log.Fatal("GetDatabase()", err)
	}

	conn, err := tst.DbConnect("./sqlite.db")
	if err != nil {
		log.Fatal("DbConnect", err)
	}

*/
