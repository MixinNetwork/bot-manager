package externals

import (
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/astaxie/beego"
	"github.com/liuzemei/bot-manager/durable"
	"github.com/liuzemei/bot-manager/session"
	"strings"
)

var (
	clientId     string
	clientSecret string
	sessionId    string
	privateKey   string
	ctx          = durable.Ctx
)

type User struct {
	UserId         string `json:"user_id"`
	FullName       string `json:"full_name"`
	IdentityNumber string `json:"identity_number"`
	AvatarURL      string `json:"avatar_url"`
	AccessToken    string `json:"access_token"`
	CreatedAt      string `json:"created_at"`
}

func init() {
	clientId = beego.AppConfig.String("clientId")
	clientSecret = beego.AppConfig.String("clientSecret")
	sessionId = beego.AppConfig.String("sessionId")
	privateKey = beego.AppConfig.String("privateKey")
}

func GetUserByCode(code string) (*bot.User, string, *session.Error) {
	token, scope, _, err := bot.OAuthGetAccessToken(ctx, clientId, clientSecret, code, "", "")
	if err != nil {
		return nil, "", parseError(err.(bot.Error))
	}
	if !strings.Contains(scope, "PROFILE:READ") {
		return nil, "", session.ForbiddenError()
	}
	user, err := bot.UserMe(ctx, token)
	if err != nil {
		return nil, "", parseError(err.(bot.Error))
	}
	return user, token, nil
}

func GetUserById(userId string) (*User, error) {
	_user, err := bot.GetUser(ctx, userId, clientId, sessionId, privateKey)
	if err != nil {
		return nil, err
	}
	var user User
	user.CreatedAt = _user.CreatedAt
	user.UserId = _user.UserId
	user.AccessToken = ""
	user.AvatarURL = _user.AvatarURL
	user.FullName = _user.FullName
	user.IdentityNumber = _user.IdentityNumber
	return &user, nil
}

func HandlePrivateKey(privateKey string) string {
	return strings.Replace(strings.Replace(string(privateKey), "\\r", "\r", -1), "\\n", "\n", -1)
}
