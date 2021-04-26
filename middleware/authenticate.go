package middleware

import (
	"regexp"
	"strings"

	"github.com/MixinNetwork/bot-manager/session"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
)

var whitelist = [][2]string{
	{"GET", "^/$"},
	{"GET", "^/_hc$"},
	{"GET", "^/users"},
	{"GET", "^/v1/user/login$"},
	{"GET", "^/config$"},
	{"GET", "^/amount$"},
	{"POST", "^/auth$"},
}

type customClaims struct {
	UserId string `json:"foo"`
	jwt.StandardClaims
}

var key = []byte(beego.AppConfig.String("claimKey"))

func authenticate() {
	var FilterUser = func(ctx *context.Context) {
		if ctx.Input.Method() == "OPTIONS" {
			return
		}
		token := ctx.Input.Header("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			handleUnauthorized(ctx)
			return
		}
		userId, err := Parse(token[7:])
		if err != nil || userId == "" {
			handleUnauthorized(ctx)
			return
		}
		ctx.Input.SetData("UserId", userId)
	}
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
}

func handleUnauthorized(ctx *context.Context) {
	for _, pp := range whitelist {
		if pp[0] != ctx.Input.Method() {
			continue
		}
		if matched, err := regexp.MatchString(pp[1], strings.ToLower(ctx.Request.URL.Path)); err == nil && matched {
			return
		}
	}
	err := session.AuthorizationError()
	session.HandleError(ctx, err)
}

func Claims(str string) (string, error) {
	claims := customClaims{str, jwt.StandardClaims{}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		return "", nil
	}
	return ss, nil
}

func Parse(ss string) (string, error) {
	token, err := jwt.ParseWithClaims(ss, &customClaims{}, func(_ *jwt.Token) (i interface{}, err error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		return claims.UserId, nil
	}
	return "", err
}
