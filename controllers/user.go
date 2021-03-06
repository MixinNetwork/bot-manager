package controllers

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/MixinNetwork/bot-manager/externals"
	"github.com/MixinNetwork/bot-manager/middleware"
	"github.com/MixinNetwork/bot-manager/models"
	"github.com/MixinNetwork/bot-manager/session"
	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

func (c *UserController) Get() {
	clientId := c.GetString("client_id")
	userId := c.Ctx.Input.GetData("UserId")
	if !checkBotManager(userId.(string), clientId, c.Ctx) {
		return
	}

	status := c.GetString("status")
	userList := models.GetUsersByClientId(clientId, status)
	c.Data["json"] = Resp{Data: userList}
	c.ServeJSON()
}

func (c *UserController) Post() {
	var user externals.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		log.Println("/controllers/user.Post Unmarshal", err)
	}
	botUser, err := externals.GetUserById(user.UserId)
	if err != nil {
		log.Println("/controllers/user.Post GetUserById", err)
	}
	c.Data["json"] = Resp{Data: botUser}
	//models.AddUser()
	c.ServeJSON()
}

func (c *UserController) Put() {
	userId := c.Ctx.Input.GetData("UserId")
	reqModel := new(struct {
		UserId   string `json:"user_id"`
		Status   string `json:"status"`
		ClientId string `json:"client_id"`
	})
	err := json.Unmarshal(c.Ctx.Input.RequestBody, reqModel)
	if err != nil {
		log.Println("/controllers/user.Put Unmarshal", err)
	}
	if !checkBotManager(userId.(string), reqModel.ClientId, c.Ctx) {
		return
	}
	models.UpdateBotUserStatus(reqModel.ClientId, reqModel.UserId, reqModel.Status)
	c.Data["json"] = Resp{Data: "ok"}
	c.ServeJSON()
}

func (c *UserController) Login() {
	code := c.GetString("code")
	user, token, err := externals.GetUserByCode(code)

	if err != nil || user == nil || token == "" {
		session.HandleError(c.Ctx, err)
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
		Data: models.User{
			AccessToken:    authToken,
			AvatarURL:      user.AvatarURL,
			FullName:       user.FullName,
			UserId:         user.UserId,
			IdentityNumber: user.IdentityNumber,
		},
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

func GetMessageUserAutoUpdate(userId, clientId string) (*models.User, error) {
	botUser := models.GetBotUser(userId, clientId)
	if botUser.IdentityNumber == "" {
		usersUser := models.GetUserById(userId)
		if usersUser.IdentityNumber == "" {
			userInfo, err := externals.GetUserById(userId)
			if err != nil {
				if strings.Contains(err.Error(), "Server Error") {
					return GetMessageUserAutoUpdate(userId, clientId)
				}
				return nil, err
			}
			modelUser := models.User{
				UserId:         userInfo.UserId,
				FullName:       userInfo.FullName,
				IdentityNumber: userInfo.IdentityNumber,
				AvatarURL:      userInfo.AvatarURL,
				AccessToken:    "",
				CreatedAt:      userInfo.CreatedAt,
			}
			models.AddUser(&modelUser)
			models.AddBotUser(&modelUser, clientId)
			botUser.FullName = userInfo.FullName
			botUser.IdentityNumber = userInfo.IdentityNumber
			botUser.AvatarURL = userInfo.AvatarURL
		} else {
			botUser.FullName = usersUser.FullName
			botUser.IdentityNumber = usersUser.IdentityNumber
			botUser.AvatarURL = usersUser.AvatarURL
			models.AddBotUser(&models.User{
				UserId:         userId,
				FullName:       usersUser.FullName,
				IdentityNumber: usersUser.IdentityNumber,
				AvatarURL:      usersUser.AvatarURL,
				AccessToken:    "",
			}, clientId)
		}
	}
	return &models.User{
		FullName:       botUser.FullName,
		IdentityNumber: botUser.IdentityNumber,
		AvatarURL:      botUser.AvatarURL,
		UserId:         userId,
	}, nil
}
