package account

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kinggigo/secret/server/config"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type Log struct {
	Token string
}

func Login(e echo.Context) error {
	//jwt.New([]byte(config.JWT_KEY))
	tokens := jwt.New(jwt.SigningMethodHS256)
	t, err := tokens.SignedString([]byte(config.JWT_KEY))
	if err != nil {
		log.Printf("err", err)
		return err
	}
	return e.JSON(200, Log{t})
}

func Happyy(e echo.Context) error {

	return e.String(200, "로그인되었다!!!")
}
