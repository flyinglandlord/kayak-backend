package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type UserInfoResponse struct {
	UserId     int       `json:"user_id"`
	UserName   string    `json:"user_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	AvatarPath string    `json:"avatar_path"`
	CreateAt   time.Time `json:"create_at"`
}

// GetUserInfo godoc
// @Schemes http
// @Description 获取用户信息
// @Success 200 {object} UserInfoResponse "用户信息"
// @Failure default {string} string "服务器错误"
// @Router /user/info [get]
// @Security ApiKeyAuth
func GetUserInfo(c *gin.Context) {
	user := model.User{}
	sqlString := `SELECT name, email, phone, avatar_url, created_at FROM "user" WHERE id = $1`
	if err := global.Database.Get(&user, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	userInfo := UserInfoResponse{
		UserId:     c.GetInt("UserId"),
		UserName:   user.Name,
		Email:      *user.Email,
		Phone:      *user.Phone,
		AvatarPath: user.AvatarURL,
		CreateAt:   user.CreatedAt,
	}
	c.JSON(200, userInfo)
}
