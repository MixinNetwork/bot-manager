package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/liuzemei/bot-manager/db"
	"github.com/liuzemei/bot-manager/externals"
	"github.com/liuzemei/bot-manager/middleware"
	"github.com/liuzemei/bot-manager/models"
	"github.com/liuzemei/bot-manager/session"
	"log"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]

func (u *UserController) Post() {
	var user externals.User
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if err != nil {
		log.Println("/controllers/user.Post Unmarshal", err)
	}
	botUser, err := externals.GetUserById(user.UserId)
	if err != nil {
		log.Println("/controllers/user.Post GetUserById", err)
	}
	u.Data["json"] = botUser
	//models.AddUser()
	u.ServeJSON()
}

type LoginRespData struct {
	AccessToken    string `json:"access_token"`
	AvatarUrl      string `json:"avatar_url"`
	FullName       string `json:"full_name"`
	UserId         string `json:"user_id"`
	IdentityNumber string `json:"identity_number"`
}

func (u *UserController) Login() {
	code := u.GetString("code")
	user, token, err := externals.GetUserByCode(code)

	if err != nil || user == nil || token == "" {
		session.HandleError(u.Ctx, err)
		return
	}

	var modelUser = models.User{
		UserId:         user.UserId,
		FullName:       user.FullName,
		IdentityNumber: user.IdentityNumber,
		AvatarURL:      user.AvatarURL,
		AccessToken:    token,
		CreatedAt:      user.CreatedAt,
	}

	models.AddUser(&modelUser)
	authToken, _ := middleware.Claims(user.UserId)

	var resp = Resp{
		Data: LoginRespData{
			AccessToken:    authToken,
			AvatarUrl:      user.AvatarURL,
			FullName:       user.FullName,
			UserId:         user.UserId,
			IdentityNumber: user.IdentityNumber,
		},
	}

	u.Data["json"] = resp
	u.ServeJSON()
}

func AddMessageUser(userId string) (*models.UserBase, error) {
	var hasUser models.User
	db.Conn.First(&hasUser, "user_id=?", userId)
	if hasUser.IdentityNumber == "" {
		userInfo, err := externals.GetUserById(userId)
		if err != nil {
			log.Println("获取userInfo出错了", err)
			return nil, err
		}
		hasUser = models.User{
			UserId:         userInfo.UserId,
			FullName:       userInfo.FullName,
			IdentityNumber: userInfo.IdentityNumber,
			AvatarURL:      userInfo.AvatarURL,
			AccessToken:    "",
			CreatedAt:      userInfo.CreatedAt,
		}
		db.Conn.Create(&hasUser)
	}
	return &models.UserBase{
		FullName:       hasUser.FullName,
		IdentityNumber: hasUser.IdentityNumber,
		AvatarURL:      hasUser.AvatarURL,
	}, nil
}
