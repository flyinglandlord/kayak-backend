package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"kayak-backend/utils"
	"net/http"
	"time"
)

type LoginInfo struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterInfo struct {
	Name     string  `json:"name" bind:"required,max=20"`
	Password string  `json:"password" bind:"required,min=6,max=20"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	VCode    string  `json:"v_code"`
}

type RegisterResponse struct {
	OldPassword string `json:"old_password" bind:"required"`
	NewPassword string `json:"new_password" bind:"required,min=6,max=20"`
}

// Login godoc
// @Schemes http
// @Description 用户登录
// @Tags auth
// @Param info body LoginInfo true "用户登陆信息"
// @Success 200 {object} LoginResponse "用户登陆反馈"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /login [post]
func Login(c *gin.Context) {
	loginRequest := LoginInfo{}
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	userInfo := model.User{}
	sqlString := `SELECT id, password FROM "user" WHERE name = $1`
	if err := global.Database.Get(&userInfo, sqlString, loginRequest.UserName); err != nil {
		c.String(http.StatusBadRequest, "用户名或密码错误")
		return
	}
	if !utils.VerifyPassword(userInfo.Password, loginRequest.Password) {
		c.String(http.StatusBadRequest, "用户名或密码错误")
		return
	}
	token, err := global.CreateSession(c, &global.Session{
		Role:   global.USER,
		UserId: userInfo.ID,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
	c.Set("Role", global.USER)
	c.Set("UserId", userInfo.ID)
}

// Register godoc
// @Schemes http
// @Description 用户注册
// @Tags auth
// @Param info body RegisterInfo true "用户注册信息"
// @Success 200 {string} string "注册成功"
// @Failure 400 {string} string "请求解析失败"/"验证码已过期"/"验证码错误"
// @Failure 409 {string} string "用户名已存在"
// @Failure default {string} string "服务器错误"
// @Router /register [post]
func Register(c *gin.Context) {
	registerRequest := RegisterInfo{}
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	userInfo := model.User{}
	sqlString := `SELECT id FROM "user" WHERE name = $1`
	if err := global.Database.Get(&userInfo, sqlString, registerRequest.Name); err == nil {
		c.String(409, "用户名已存在")
		return
	}
	rawCode := global.Redis.Get(c, registerRequest.Email)
	if rawCode.Err() != nil {
		c.String(http.StatusBadRequest, "验证码已过期")
		return
	} else if rawCode.Val() != registerRequest.VCode {
		c.String(http.StatusBadRequest, "验证码错误")
		return
	} else {
		global.Redis.Del(c, registerRequest.Email)
	}
	userInfo.Name = registerRequest.Name
	var err error
	userInfo.Password, err = utils.EncryptPassword(registerRequest.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	userInfo.Email = &registerRequest.Email
	userInfo.Phone = registerRequest.Phone
	sqlString = `INSERT INTO "user" (name, password, email, phone, created_at, nick_name) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	if err := global.Database.Get(&userInfo.ID, sqlString, userInfo.Name, userInfo.Password,
		userInfo.Email, userInfo.Phone, time.Now().Local(), userInfo.Name); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "注册成功")
}

// ResetPassword godoc
// @Schemes http
// @Description 用户修改密码
// @Tags auth
// @Param info body RegisterResponse true "用户修改密码信息"
// @Success 200 {string} string "修改成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 409 {string} string "用户不存在"
// @Failure default {string} string "服务器错误"
// @Router /reset-password [post]
// @Security ApiKeyAuth
func ResetPassword(c *gin.Context) {
	registerRequest := RegisterResponse{}
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	userId := c.GetInt("UserId")
	userInfo := model.User{}
	sqlString := `SELECT id, password FROM "user" WHERE id = $1`
	if err := global.Database.Get(&userInfo, sqlString, userId); err != nil {
		c.String(409, "用户不存在")
		return
	}
	if !utils.VerifyPassword(userInfo.Password, registerRequest.OldPassword) {
		c.String(http.StatusBadRequest, "旧密码错误")
		return
	}
	var err error
	userInfo.Password, err = utils.EncryptPassword(registerRequest.NewPassword)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `UPDATE "user" SET password = $1 WHERE id = $2`
	if _, err := global.Database.Exec(sqlString, userInfo.Password, userInfo.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "修改成功")
}

// Logout godoc
// @Schemes http
// @Description 用户退出
// @Tags auth
// @Success 200 {string} string "退出成功"
// @Failure default {string} string "服务器错误"
// @Router /logout [get]
// @Security ApiKeyAuth
func Logout(c *gin.Context) {
	err := global.DeleteSession(c, c.Request.Header.Get(global.TokenHeader))
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
	}
	c.String(http.StatusOK, "退出成功")
}

// SendEmail godoc
// @Schemes http
// @Description 发送邮件
// @Tags auth
// @Param email query string true "邮箱"
// @Success 200 {string} string "发送成功"
// @Failure 400 {string} string "验证码存储失败"
// @Failure default {string} string "服务器错误"
// @Router /send-email [post]
func SendEmail(c *gin.Context) {
	fmt.Println(c.Query("email"))
	vCode, err := utils.SendEmailValidate(c.Query("email"))
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	err = global.Redis.Set(c, c.Query("email"), vCode, time.Minute*5).Err()
	if err != nil {
		c.String(http.StatusInternalServerError, "验证码存储失败")
		return
	}
	c.String(http.StatusOK, "发送成功")
}
