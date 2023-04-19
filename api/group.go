package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type GroupResponse struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserId      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}
type GroupCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateGroup godoc
// @Schemes http
// @Description 创建小组
// @Tags group
// @Param group body GroupCreateRequest true "小组信息"
// @Success 200 {object} GroupResponse "小组信息"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /group/create [post]
// @Security ApiKeyAuth
func CreateGroup(c *gin.Context) {
	var request GroupCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `INSERT INTO "group" (name, description, user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var groupId int
	if err := global.Database.Get(&groupId, sqlString, request.Name, request.Description, c.GetInt("UserId"), time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var group model.Group
	sqlString = `SELECT * FROM "group" WHERE id = $1`
	if err := global.Database.Get(&group, sqlString, groupId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, GroupResponse{
		Id:          group.Id,
		Name:        group.Name,
		Description: group.Description,
		UserId:      group.UserId,
		CreatedAt:   group.CreatedAt,
	})
}

// DeleteGroup godoc
// @Schemes http
// @Description 删除小组
// @Tags group
// @Param id path int true "小组ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "小组不存在"
// @Failure default {string} string "服务器错误"
// @Router /group/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteGroup(c *gin.Context) {
	sqlString := `SELECT user_id FROM "group" WHERE id = $1`
	var groupUserId int
	if err := global.Database.Get(&groupUserId, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "小组不存在")
		return
	}
	if role, _ := c.Get("Role"); groupUserId != c.GetInt("UserId") && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM "group" WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

type AllUserResponse struct {
	TotalCount int                `json:"total_count"`
	User       []UserInfoResponse `json:"user"`
}

// GetUsersInGroup godoc
// @Schemes http
// @Description 获取小组成员
// @Tags group
// @Param id path int true "小组ID"
// @Success 200 {object} AllUserResponse "用户信息"
// @Failure 404 {string} string "小组不存在"
// @Failure default {string} string "服务器错误"
// @Router /group/all_user/{id} [get]
// @Security ApiKeyAuth
func GetUsersInGroup(c *gin.Context) {
	sqlString := `SELECT user_id FROM "group" WHERE id = $1`
	var groupUserId int
	if err := global.Database.Get(&groupUserId, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "小组不存在")
		return
	}
	var users []model.User
	sqlString = `SELECT * FROM "user" WHERE id IN (SELECT user_id FROM group_member WHERE group_id = $1)`
	if err := global.Database.Select(&users, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var userResponses []UserInfoResponse
	for _, user := range users {
		userResponses = append(userResponses, UserInfoResponse{
			UserId:     user.ID,
			UserName:   user.Name,
			Email:      user.Email,
			Phone:      user.Phone,
			AvatarPath: user.AvatarURL,
			CreateAt:   user.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, AllUserResponse{
		TotalCount: len(userResponses),
		User:       userResponses,
	})
}

// AddUserToGroup godoc
// @Schemes http
// @Description 添加用户到小组
// @Tags group
// @Param id path int true "小组ID"
// @Param user_id query int true "用户ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "小组不存在"/"用户不存在"
// @Failure default {string} string "服务器错误"
// @Router /group/add/{id} [post]
// @Security ApiKeyAuth
func AddUserToGroup(c *gin.Context) {
	sqlString := `SELECT user_id FROM "group" WHERE id = $1`
	var groupUserId int
	if err := global.Database.Get(&groupUserId, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "小组不存在")
		return
	}
	sqlString = `SELECT id FROM "user" WHERE id = $1`
	var userId int
	if err := global.Database.Get(&userId, sqlString, c.Query("user_id")); err != nil {
		c.String(http.StatusNotFound, "用户不存在")
		return
	}
	if role, _ := c.Get("Role"); groupUserId != c.GetInt("UserId") && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO group_member (user_id, group_id) VALUES ($1, $2)`
	if _, err := global.Database.Exec(sqlString, c.Query("user_id"), c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveUserFromGroup godoc
// @Schemes http
// @Description 从小组移除用户
// @Tags group
// @Param id path int true "小组ID"
// @Param user_id query int true "用户ID"
// @Success 200 {string} string "移除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "小组不存在"
// @Failure default {string} string "服务器错误"
// @Router /group/remove/{id} [delete]
// @Security ApiKeyAuth
func RemoveUserFromGroup(c *gin.Context) {
	sqlString := `SELECT user_id FROM "group" WHERE id = $1`
	var groupUserId int
	if err := global.Database.Get(&groupUserId, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "小组不存在")
		return
	}
	if role, _ := c.Get("Role"); groupUserId != c.GetInt("UserId") && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM group_member WHERE user_id = $1 AND group_id = $2`
	if _, err := global.Database.Exec(sqlString, c.Query("user_id"), c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "移除成功")
}
