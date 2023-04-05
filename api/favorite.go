package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"net/http"
	"time"
)

type FavoriteProblemResponse struct {
	ProblemID     int `json:"problem_id"`
	ProblemTypeID int `json:"problem_type_id"`
}

// AddProblemToFavorite godoc
// @Schemes http
// @Description 添加题目到收藏夹（只有管理员和题目创建者能添加私有题目）
// @Param id path int true "题目ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/favorite/{id} [post]
func AddProblemToFavorite(c *gin.Context) {
	problemId := c.Param("id")
	userId := c.GetInt("UserId")
	var problemUserId int
	var isPublic bool
	sqlString := "SELECT is_public, user_id FROM problem_type WHERE id = $1"
	if err := global.Database.Get(&struct {
		IsPublic bool `db:"is_public"`
		UserId   int  `db:"user_id"`
	}{IsPublic: isPublic, UserId: problemUserId}, sqlString, problemId); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && problemUserId != userId && !isPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = "INSERT INTO user_favorite_problem (user_id, problem_id, created_at) VALUES ($1, $2, $3)"
	if _, err := global.Database.Exec(sqlString, userId, problemId, time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveProblemFromFavorite godoc
// @Schemes http
// @Description 从收藏夹移除题目
// @Param id path int true "题目ID"
// @Success 200 {string} string "移除成功"
// @Failure default {string} string "服务器错误"
// @Router /problem/unfavorite/{id} [post]
func RemoveProblemFromFavorite(c *gin.Context) {
	problemId := c.Param("id")
	userId := c.GetInt("UserId")
	sqlString := "DELETE FROM user_favorite_problem WHERE user_id = $1 AND problem_id = $2"
	if _, err := global.Database.Exec(sqlString, userId, problemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "移除成功")
}

// AddProblemsetToFavorite godoc
// @Schemes http
// @Description 添加题集到收藏夹（只有管理员和题集创建者能添加私有题集）
// @Param id path int true "题集ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /problemset/favorite/{id} [post]
func AddProblemsetToFavorite(c *gin.Context) {
	problemsetId := c.Param("id")
	userId := c.GetInt("UserId")
	var problemsetUserId int
	var isPublic bool
	sqlString := "SELECT is_public, user_id FROM problemset WHERE id = $1"
	if err := global.Database.Get(&struct {
		IsPublic bool `db:"is_public"`
		UserId   int  `db:"user_id"`
	}{IsPublic: isPublic, UserId: problemsetUserId}, sqlString, problemsetId); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && problemsetUserId != userId && !isPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = "INSERT INTO user_favorite_problemset (user_id, problemset_id, created_at) VALUES ($1, $2, $3)"
	if _, err := global.Database.Exec(sqlString, userId, problemsetId, time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveProblemsetFromFavorite godoc
// @Schemes http
// @Description 从收藏夹移除题集
// @Param id path int true "题集ID"
// @Success 200 {string} string "移除成功"
// @Failure default {string} string "服务器错误"
// @Router /problemset/unfavorite/{id} [post]
func RemoveProblemsetFromFavorite(c *gin.Context) {
	problemsetId := c.Param("id")
	userId := c.GetInt("UserId")
	sqlString := "DELETE FROM user_favorite_problemset WHERE user_id = $1 AND problemset_id = $2"
	if _, err := global.Database.Exec(sqlString, userId, problemsetId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "移除成功")
}
