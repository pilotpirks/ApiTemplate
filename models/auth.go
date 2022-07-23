package models

import (
	"image/color"

	"github.com/golang-jwt/jwt"
	"github.com/mojocn/base64Captcha"
)

var capthaStore = base64Captcha.DefaultMemStore

var capthaConfig = base64Captcha.DriverString{
	Height:          60,
	Width:           180,
	ShowLineOptions: 0,
	NoiseCount:      5,
	Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
	Length:          5,
	Fonts: []string{
		"3Dumb.ttf",
		// "ApothecaryFont.ttf",
		// "Comismsh.ttf",
		// "DENNEthree-dee.ttf",
		// "DeborahFancyDress.ttf",
		// "Flim-Flam.ttf",
		// "RitaSmith.ttf",
		// "actionj.ttf",
		// "chromohv.ttf",
		// "wqy-microhei.ttc",
	},
	BgColor: &color.RGBA{0, 0, 0, 0},
}

type Auth struct {
	CaptchaId      string `json:"captcha_id" validate:"required"`
	CaptchaCode    string `json:"captcha_code" validate:"required"`
	Name           string `json:"name" validate:"required"`
	Password       string `json:"password" validate:"required"`
	DoublePassword string `json:"double_password"`
}

type JwtCustomClaims struct {
	Name    string
	IsAdmin bool
	jwt.StandardClaims
}

// Create token with claims
func (claims *JwtCustomClaims) CreateToken() *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

type Captcha struct {
	Id  string `json:"id"`
	Img string `json:"image"`
	Err string `json:"error"`
}

func (c Captcha) Generate() Captcha {
	driver := capthaConfig.ConvertFonts()
	newCaptcha := base64Captcha.NewCaptcha(driver, capthaStore)
	id, b64s, err := newCaptcha.Generate()
	if err != nil {
		c.Err = err.Error()
		return c
	}

	c.Id = id
	c.Img = b64s
	return c
}

func (c Captcha) Verify(id, code string) bool {
	return capthaStore.Verify(id, code, true)
}
