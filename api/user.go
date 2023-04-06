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
	Email      *string   `json:"email"`
	Phone      *string   `json:"phone"`
	AvatarPath string    `json:"avatar_path"`
	CreateAt   time.Time `json:"create_at"`
}

// GetUserInfoById godoc
// @Schemes http
// @Description 根据ID获取用户信息
// @Param user_id path int true "用户ID"
// @Success 200 {object} UserInfoResponse "用户信息"
// @Failure 404 {string} string "用户不存在"
// @Failure default {string} string "服务器错误"
// @Router /user/info/{user_id} [get]
// @Security ApiKeyAuth
func GetUserInfoById(c *gin.Context) {
	user := model.User{}
	sqlString := `SELECT id, name, email, phone, avatar_url, created_at FROM "user" WHERE id = $1`
	if err := global.Database.Get(&user, sqlString, c.Param("user_id")); err != nil {
		c.String(http.StatusNotFound, "用户不存在")
		return
	}
	userInfo := UserInfoResponse{
		UserId:     user.ID,
		UserName:   user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		AvatarPath: user.AvatarURL,
		CreateAt:   user.CreatedAt,
	}
	c.JSON(http.StatusOK, userInfo)
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
		Email:      user.Email,
		Phone:      user.Phone,
		AvatarPath: user.AvatarURL,
		CreateAt:   user.CreatedAt,
	}
	c.JSON(http.StatusOK, userInfo)
}

type UserInfoRequest struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

// UpdateUserInfo godoc
// @Schemes http
// @Description 更新用户信息
// @Param info body UserInfoRequest true "用户信息"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /user/update [put]
// @Security ApiKeyAuth
func UpdateUserInfo(c *gin.Context) {
	user := UserInfoRequest{}
	sqlString := `UPDATE "user" SET name = $1, email = $2, phone = $3 WHERE id = $4`
	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, "请求格式错误")
		return
	}
	if _, err := global.Database.Exec(sqlString, user.Name, user.Email, user.Phone, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "更新成功")
}

type UserNotesResponse struct {
	TotalCount int            `json:"total_count"`
	Notes      []NoteResponse `json:"notes"`
}

// GetUserWrongRecords godoc
// @Schemes http
// @Description 获取当前登录用户的所有错题记录
// @Success 200 {object} AllWrongRecordResponse "错题记录列表"
// @Failure default {string} string "服务器错误"
// @Router /user/wrong_record [get]
// @Security ApiKeyAuth
func GetUserWrongRecords(c *gin.Context) {
	var wrongRecords []model.WrongRecord
	sqlString := `SELECT * FROM user_wrong_record WHERE user_id = $1`
	if err := global.Database.Select(&wrongRecords, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var wrongRecordResponses []WrongRecordResponse
	for _, wrongRecord := range wrongRecords {
		wrongRecordResponses = append(wrongRecordResponses, WrongRecordResponse{
			ProblemId: wrongRecord.ProblemId,
			Count:     wrongRecord.Count,
			CreatedAt: wrongRecord.CreatedAt,
			UpdatedAt: wrongRecord.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, AllWrongRecordResponse{
		TotalCount: len(wrongRecordResponses),
		Records:    wrongRecordResponses,
	})
}

