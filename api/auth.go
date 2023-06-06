package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"kayak-backend/global"
	"kayak-backend/model"
	"kayak-backend/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	lastSentTimesKey = "last_sent_times"
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

type ResetPasswordInfo struct {
	UserName    string `json:"username" bind:"required"`
	VerifyCode  string `json:"verify_code" bind:"required"`
	NewPassword string `json:"new_password" bind:"required,min=6,max=20"`
}

// Login godoc
// @Schemes http
// @Description 用户登录
// @Tags Authentication
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
	if loginRequest.UserName == "" || loginRequest.Password == "" {
		c.String(http.StatusBadRequest, "用户名或密码错误")
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
// @Tags Authentication
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
	userInfo.Email = registerRequest.Email
	userInfo.Phone = registerRequest.Phone
	sqlString = `INSERT INTO "user" (name, password, email, phone, created_at, nick_name) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	if err := global.Database.Get(&userInfo.ID, sqlString, userInfo.Name, userInfo.Password,
		userInfo.Email, userInfo.Phone, time.Now().Local(), userInfo.Name); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "注册成功")
}

// ChangePassword godoc
// @Schemes http
// @Description 用户修改密码
// @Tags Authentication
// @Param info body RegisterResponse true "用户修改密码信息"
// @Success 200 {string} string "修改成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 409 {string} string "用户不存在"
// @Failure default {string} string "服务器错误"
// @Router /change-password [post]
// @Security ApiKeyAuth
func ChangePassword(c *gin.Context) {
	registerRequest := RegisterResponse{}
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	userId := c.GetInt("UserId")
	userInfo := model.User{}
	sqlString := `SELECT id, password FROM "user" WHERE id = $1`
	if err := global.Database.Get(&userInfo, sqlString, userId); err != nil {
		c.String(409, "修改密码失败")
		return
	}
	if !utils.VerifyPassword(userInfo.Password, registerRequest.OldPassword) {
		c.String(http.StatusBadRequest, "修改密码失败")
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

// ResetPassword godoc
// @Schemes http
// @Description 用户重置密码，需要先向邮箱发送一封邮件
// @Tags Authentication
// @Param info body ResetPasswordInfo true "用户修改密码信息"
// @Success 200 {string} string "修改成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 409 {string} string "用户不存在"
// @Failure default {string} string "服务器错误"
// @Router /reset-password [post]
// @Security ApiKeyAuth
func ResetPassword(c *gin.Context) {
	// 验证Redis内的邮箱验证码是否正确，然后修改密码
	resetPasswordInfo := ResetPasswordInfo{}
	if err := c.ShouldBindJSON(&resetPasswordInfo); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	userInfo := model.User{}
	sqlString := `SELECT * FROM "user" WHERE name = $1`
	if err := global.Database.Get(&userInfo, sqlString, resetPasswordInfo.UserName); err != nil {
		c.String(409, "修改密码失败")
		return
	}
	rawCode := global.Redis.Get(c, userInfo.Email)
	if rawCode.Err() != nil {
		c.String(http.StatusBadRequest, "修改密码失败")
		return
	} else if rawCode.Val() != resetPasswordInfo.VerifyCode {
		c.String(http.StatusBadRequest, "修改密码失败")
		return
	} else {
		global.Redis.Del(c, userInfo.Email)
	}
	var err error
	userInfo.Password, err = utils.EncryptPassword(resetPasswordInfo.NewPassword)
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
// @Tags Authentication
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
// @Tags Authentication
// @Param email query string true "邮箱"
// @Success 200 {string} string "发送成功"
// @Failure 400 {string} string "验证码存储失败"
// @Failure default {string} string "服务器错误"
// @Router /send-email [post]
func SendEmail(c *gin.Context) {
	// fmt.Println(c.Query("email"))
	email := c.Query("email")
	currentTime := time.Now()
	lastSentTimeStr, err := global.Redis.ZScore(c, lastSentTimesKey, email).Result()
	if err == redis.Nil || (err == nil && currentTime.Sub(time.Unix(int64(lastSentTimeStr), 0)) >= time.Minute) {
		vCode, err := utils.SendEmailValidate(c.Query("email"))
		if err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		err = global.Redis.ZAdd(c, lastSentTimesKey, &redis.Z{
			Score:  float64(currentTime.Unix()),
			Member: email,
		}).Err()
		if err != nil {
			c.String(http.StatusInternalServerError, "验证码存储失败")
			return
		}
		err = global.Redis.Set(c, email, vCode, time.Minute*5).Err()
		if err != nil {
			c.String(http.StatusInternalServerError, "验证码存储失败")
			return
		}
		c.String(http.StatusOK, "发送成功")
	} else {
		c.String(http.StatusBadRequest, "发送过于频繁")
	}
}

type WeixinLoginInfo struct {
	Code string `json:"code"`
}

type WeixinReturnInfo struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type WeixinLoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

var (
	url = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

// WeixinLogin godoc
// @Schemes http
// @Description 微信登录
// @Tags Authentication
// @Param code body WeixinLoginInfo true "微信登录信息"
// @Success 200 {object} WeixinLoginResponse "用户登陆反馈"
// @Failure 400 {string} string "请求解析失败"
// @Failure 409 {string} string "用户不存在"
// @Failure default {string} string "服务器错误"
// @Router /weixin-login [post]
func WeixinLogin(c *gin.Context) {
	weixinLoginInfo := WeixinLoginInfo{}
	if err := c.ShouldBindJSON(&weixinLoginInfo); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	// 获取微信用户信息
	var weixinReturnInfo WeixinReturnInfo
	httpState, bytes := utils.Get(fmt.Sprintf(url, global.AppID, global.AppSecret, weixinLoginInfo.Code))
	if httpState != 200 {
		_ = fmt.Errorf("获取SessionKey失败, Http状态码: %d", httpState)
		c.String(http.StatusBadRequest, weixinReturnInfo.ErrMsg+" 获取SessionKey失败")
	}
	e := json.Unmarshal(bytes, &weixinReturnInfo)
	if e != nil {
		_ = fmt.Errorf("json解析失败")
		c.String(http.StatusBadRequest, weixinReturnInfo.ErrMsg+" json解析失败")
	}

	// 查询用户是否存在
	userInfo := model.User{}
	sqlString := `SELECT * FROM "user" WHERE open_id = $1`
	if err := global.Database.Get(&userInfo, sqlString, weixinReturnInfo.OpenID); err != nil {
		// 添加一个用户
		sqlString = `INSERT INTO "user" (open_id, name, email, phone, password, created_at, nick_name) VALUES ($1, '', '', '', '', now(), $2) RETURNING id`
		if err := global.Database.Get(&userInfo, sqlString, weixinReturnInfo.OpenID, "新注册用户"+strconv.FormatInt(rand.Int63(), 10)); err != nil {
			c.String(http.StatusInternalServerError, "新注册用户添加失败")
			return
		}
		token, err := global.CreateSession(c, &global.Session{
			Role:   global.USER,
			UserId: userInfo.ID,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "新注册用户Token生成失败")
			return
		}
		c.Set("Role", global.USER)
		c.Set("UserId", userInfo.ID)
		c.JSON(http.StatusOK, WeixinLoginResponse{
			Code:    201,
			Message: "待完善信息",
			Token:   token,
		})
		return
	}
	token, err := global.CreateSession(c, &global.Session{
		Role:   global.USER,
		UserId: userInfo.ID,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "已注册用户Token生成失败")
		return
	}
	c.Set("Role", global.USER)
	c.Set("UserId", userInfo.ID)
	c.JSON(http.StatusOK, WeixinLoginResponse{
		Code:    200,
		Message: "登录成功",
		Token:   token,
	})
}

// WeixinBind godoc
// @Schemes http
// @Description 微信绑定
// @Tags Authentication
// @Param code body WeixinLoginInfo true "微信登录信息"
// @Success 200 {string} string "绑定成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /weixin-bind [post]
// @Security ApiKeyAuth
func WeixinBind(c *gin.Context) {
	weixinLoginInfo := WeixinLoginInfo{}
	if err := c.ShouldBindJSON(&weixinLoginInfo); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	// 获取微信用户信息
	var weixinReturnInfo WeixinReturnInfo
	httpState, bytes := utils.Get(fmt.Sprintf(url, global.AppID, global.AppSecret, weixinLoginInfo.Code))
	if httpState != 200 {
		_ = fmt.Errorf("获取SessionKey失败, Http状态码: %d", httpState)
		c.String(http.StatusBadRequest, weixinReturnInfo.ErrMsg+" 获取SessionKey失败")
	}
	e := json.Unmarshal(bytes, &weixinReturnInfo)
	if e != nil {
		_ = fmt.Errorf("json解析失败")
		c.String(http.StatusBadRequest, weixinReturnInfo.ErrMsg+" json解析失败")
	}

	// 查询微信是否已被绑定
	userInfo := model.User{}
	sqlString := `SELECT * FROM "user" WHERE open_id = $1`
	if err := global.Database.Get(&userInfo, sqlString, weixinReturnInfo.OpenID); err == nil {
		c.String(http.StatusBadRequest, "该微信已被绑定")
		return
	}
	// 查询用户是否已经绑定微信
	sqlString = `SELECT * FROM "user" WHERE id = $1`
	if err := global.Database.Get(&userInfo, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusBadRequest, "用户不存在")
		return
	}
	if userInfo.OpenId != nil && *userInfo.OpenId != "" {
		c.String(http.StatusBadRequest, "该用户已绑定微信")
		return
	}
	// 绑定用户
	sqlString = `UPDATE "user" SET open_id = $1 WHERE id = $2`
	if _, err := global.Database.Exec(sqlString, weixinReturnInfo.OpenID, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "绑定成功")
}

type WeixinCompleteInfo struct {
	UserName   string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	VerifyCode string `json:"verify_code" binding:"required"`
}

// WeixinComplete godoc
// @Schemes http
// @Description 微信完善信息
// @Tags Authentication
// @Param info body WeixinCompleteInfo true "微信完善信息"
// @Success 200 {string} string "完善成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /weixin-complete [post]
// @Security ApiKeyAuth
func WeixinComplete(c *gin.Context) {
	var weixinCompleteInfo WeixinCompleteInfo
	if err := c.ShouldBindJSON(&weixinCompleteInfo); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	// 查询用户是否存在
	userInfo := model.User{}
	sqlString := `SELECT * FROM "user" WHERE id = $1`
	if err := global.Database.Get(&userInfo, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusBadRequest, "用户不存在")
		return
	}
	if userInfo.Name != "" {
		c.String(http.StatusBadRequest, "该用户已完善信息")
		return
	}
	// 验证验证码
	rawCode := global.Redis.Get(c, weixinCompleteInfo.Email)
	if rawCode.Err() != nil {
		c.String(http.StatusBadRequest, "验证码已过期")
		return
	} else if rawCode.Val() != weixinCompleteInfo.VerifyCode {
		c.String(http.StatusBadRequest, "验证码错误")
		return
	} else {
		global.Redis.Del(c, weixinCompleteInfo.Email)
	}
	// 完善用户信息
	sqlString = `UPDATE "user" SET name = $1, email = $2, password = $3 WHERE id = $4`
	encryptedPassword, err := utils.EncryptPassword(weixinCompleteInfo.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if _, err := global.Database.Exec(sqlString, weixinCompleteInfo.UserName, weixinCompleteInfo.Email, encryptedPassword, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "完善成功")
}
