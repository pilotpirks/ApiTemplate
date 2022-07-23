package application

import (
	"config"
	"dbsql"
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	instance *Application
	once     sync.Once
)

type Application struct {
	DB  *sqlx.DB
	Cfg *config.Config
}

func Get() (*Application, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	db, err := dbsql.GetDatabase(cfg.GetDBType())
	if err != nil {
		return nil, err
	}

	conn, err := db.DbConnect(cfg.GetDBConnStr())
	if err != nil {
		return nil, err
	}

	once.Do(func() {
		instance = &Application{
			DB:  conn,
			Cfg: cfg,
		}
	})

	return instance, nil
}
