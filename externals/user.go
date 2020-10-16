package externals

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/astaxie/beego"
	"github.com/liuzemei/bot-manager/durable"
	"github.com/liuzemei/bot-manager/session"
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
	privateKey = HandlePrivateKey(privateKey)
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

func GetToken(method, uri, body string) (string, error) {
	accessToken, err := bot.SignAuthenticationToken(clientId, sessionId, privateKey, method, uri, body)
	if err != nil {
		return "", err
	}
	return accessToken, err
}

func UploadFile(body io.Reader, size int64) (*bot.Attachment, error) {
	attachment, err := bot.CreateAttachment(durable.Ctx, clientId, sessionId, privateKey)
	if err != nil {
		log.Println("externals uploadFile CreateAttachment", err)
		return nil, err
	}
	uploadToAMZ(attachment.UploadUrl, body, size)

	return attachment, nil
}
func uploadToAMZ(url string, body io.Reader, size int64) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Println("uploadToAMZ err", err)
	}
	req.Header.Add("x-amz-acl", "public-read")
	req.Header.Add("Content-Type", "application/octet-stream")
	req.ContentLength = size
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("执行失败了 err", err)
		return
	}
	response.Body.Close()
}

func HandlePrivateKey(privateKey string) string {
	return strings.Replace(strings.Replace(string(privateKey), "\\r", "\r", -1), "\\n", "\n", -1)
}
