package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"models"
	"net/http"
	"strings"
	"time"
	"utils"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// ================================================== auth ==================================================

func checkAuth(c echo.Context) (models.User, error) {
	cc := c.(*CustomContext)

	var err error
	var raw []byte
	var claims models.JwtCustomClaims
	var user models.User

	tokenString := c.Request().Header.Get("Authorization")
	if len(tokenString) == 0 || tokenString == "null" {
		err = errors.New("token not found")
		return user, err
	}

	seg := strings.Split(tokenString, ".")[1]

	raw, err = jwt.DecodeSegment(seg)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(raw, &claims)
	if err != nil {
		return user, err
	}

	user, err = models.GetBy[models.User](cc.DB, "SELECT * FROM users WHERE name=$1", claims.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("запись не найдена")
		}
		return user, err
	}

	signKey := []byte(user.JwtKey)

	token, err := jwt.ParseWithClaims(tokenString, &models.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		return user, err
	}

	if !token.Valid {
		err = errors.New("checkAuth()::invalid token")
		return user, err
	}

	return user, nil
}

// e.POST("/auth/captcha", captcha)
func captcha(c echo.Context) error {
	var captcha models.Captcha

	data := struct {
		Captcha models.Captcha `json:"captcha"`
	}{captcha.Generate()}

	return c.JSON(http.StatusOK, data)
}

// e.POST("/auth/login", login)
func login(c echo.Context) error {
	cc := c.(*CustomContext)

	var err error
	var captcha models.Captcha
	var claims *models.JwtCustomClaims
	var token *jwt.Token
	var tokenStr string

	var p models.Auth
	if err = c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "Invalid Data" + err.Error()})
	}

	name := p.Name
	password := p.Password
	captchaId := p.CaptchaId
	captchaCode := p.CaptchaCode

	// get user by name
	user, err := models.GetBy[models.User](cc.DB, "SELECT * FROM users WHERE name=$1", name)
	if err != nil {
		err = errors.New("Пользователь не найден")
		goto EXIT
	}

	if name != user.Name {
		err = errors.New("Неправильное имя пользователя")
		goto EXIT
	}

	if !utils.CheckPasswordHash(user.Password, password) {
		err = errors.New("Неверный пароль")
		goto EXIT
	}

	if !captcha.Verify(captchaId, captchaCode) {
		err = errors.New("Неверный код капчи")
		goto EXIT
	}

	// Set custom claims
	claims = &models.JwtCustomClaims{
		Name:    name,  // name
		IsAdmin: false, // is admin
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		},
	}

	// Create token with claims
	token = claims.CreateToken()

	// Generate encoded token and send it as response.
	tokenStr, err = token.SignedString([]byte(user.JwtKey))
	if err != nil {
		goto EXIT
	}

EXIT:
	if err != nil {
		models.ToLog(cc.DB, &models.Log{Level: "Error", MemberId: "", Action: utils.GetFuncName(), Text: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]string{"status": err.Error()})
	}

	models.ToLog(cc.DB, &models.Log{Level: "Success", MemberId: "", Action: utils.GetFuncName(), Text: "logged"})
	data := struct {
		IdToken string `json:"id_token"`
	}{tokenStr}

	return c.JSON(http.StatusOK, data)
}

func register(c echo.Context) error {
	cc := c.(*CustomContext)

	var err error
	var captcha models.Captcha
	var bcryptPassword string

	var p models.Auth
	if err = c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "Invalid Data" + err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "Invalid Data" + err.Error()})
	}

	name := p.Name
	password := p.Password
	dPassword := p.DoublePassword
	captchaId := p.CaptchaId
	captchaCode := p.CaptchaCode

	if !captcha.Verify(captchaId, captchaCode) {
		err = errors.New("Неверный код капчи")
		goto EXIT
	}

	if len(password) == 0 || len(dPassword) == 0 {
		err = errors.New("Пароли не должны быть пустыми")
		goto EXIT
	}

	if password != dPassword {
		err = errors.New("Пароли не совпадают")
		goto EXIT
	}

	if models.IfExistExec(cc.DB, "SELECT EXISTS(SELECT 1 FROM users WHERE name=$1 LIMIT 1)", name) {
		err = errors.New("Пользователь с таким именем уже создан")
		goto EXIT
	}

	// -------------------------------------------------------------------------------------

	bcryptPassword, err = utils.HashPassword(password)
	if err != nil {
		goto EXIT
	}

	err = models.AddUser(
		&models.User{
			Id:       utils.GenUUID(),
			Name:     name,
			Password: bcryptPassword,
			JwtKey:   utils.GenUUID(),
		},
		cc.DB,
	)
	if err != nil {
		goto EXIT
	}

	// -------------------------------------------------------------------------------------

EXIT:
	if err != nil {
		models.ToLog(cc.DB, &models.Log{Level: "Error", MemberId: "", Action: utils.GetFuncName(), Text: err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]string{"status": err.Error()})
	}

	models.ToLog(cc.DB, &models.Log{Level: "Success", MemberId: "", Action: utils.GetFuncName(), Text: "registered"})
	return c.JSON(http.StatusOK, true)
}

// ================================================== pages ==================================================
// https://github.com/swaggo/swag

// @Summary      Check access
// @Description  Check service access
// @Tags         root
// @Accept       */*
// @Produce      json
// @Success 200 {boolean} boolean true
// @Failure      400
// @Router       /ping [get]
func ping(c echo.Context) error {
	return c.JSON(http.StatusOK, true)
}

// e.POST("/users/all", allUsers)
// @Summary Get list of users
// @Description Get list of users
// @Security ApiKeyAuth
// @Tags user
// @Produce json
// @Success 200 {array} models.User
// @failure 400 {string} string "error string"
// @Router /users/all [post]
func allUsers(c echo.Context) error {
	cc := c.(*CustomContext)

	users, err := models.SelectBy[models.User](cc.DB, "SELECT * FROM users")
	if err != nil {
		goto EXIT
	}

EXIT:
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// @Summary Update user parameters
// @Description Update user parameters
// @Security ApiKeyAuth
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.User true "Updated models.User"
// @Success 200 {boolean} boolean true
// @failure 400 {string} string
// @Router /users/update [post]
func updateUser(c echo.Context) error {
	cc := c.(*CustomContext)
	var err error

	var user models.User
	if err = c.Bind(&user); err != nil {
		err = errors.New("Invalid Data" + err.Error())
		goto EXIT
	}

	err = models.UpdateUser(&user, cc.DB)
	if err != nil {
		goto EXIT
	}

EXIT:
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, true)
}
