package main

import (
	"application"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type CustomContext struct {
	echo.Context
	*application.Application
}

func main() {
	// rand.Seed(time.Now().UnixNano())

	app, err := application.Get()
	if err != nil {
		log.Fatal(err)
	}

	// db.MustExec(schema)

	// ----------------------------------- init server -----------------------------------

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	/* add db, config */
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c, app}
			return next(cc)
		}
	})

	/* If you are using Gzip middleware you should add the swagger endpoint to skipper */
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	// e.File("/", "../admin_panel/dist/index.html")
	// e.Static("/", "../admin_panel/dist")
	// e.File("/favicon.ico", "../admin_panel/dist/logo.png")

	/* ---------------------------------- api docs ---------------------------------- */

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	/* ---------------------------------- register / auth ---------------------------------- */

	e.POST("/auth/captcha", captcha)
	e.POST("/auth/login", login)
	e.POST("/auth/register", register)

	/* ---------------------------------- pages ---------------------------------- */

	e.GET("/ping", ping)

	u := e.Group("/users")
	u.POST("/all", allUsers)
	u.POST("/update", updateUser)

	/* ---------------------------------- run server ---------------------------------- */

	go func() {
		if app.Cfg.Settings.Debug {
			if err := e.Start("127.0.0.1:9090"); err != nil && err != http.ErrServerClosed {
				e.Logger.Fatal("echo exit", err)
			}
		} else {
			if err := e.StartTLS("127.0.0.1::9090", "server.pem", "server.key"); err != nil && err != http.ErrServerClosed {
				e.Logger.Fatal("echo exit", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	app.DB.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}

/*
openssl genrsa -out server.key 2048
openssl req -new -x509 -key server.key -out server.pem -days 3650

go build -trimpath -ldflags="-w -s" -o ./server

kill $(pgrep server)
nohup ./server &
*/
