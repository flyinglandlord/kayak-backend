package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

// AddProblemToFavorite godoc
// @Schemes http
// @Description 添加题目到收藏夹（只有管理员和题目创建者能添加私有题目）
// @Param id path int true "题目ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/favorite/{id} [post]
// @Security ApiKeyAuth
func AddProblemToFavorite(c *gin.Context) {
	var problem model.ProblemType
	problemId := c.Param("id")
	userId := c.GetInt("UserId")
	sqlString := `SELECT * FROM problem_type WHERE id = $1`
	if err := global.Database.Get(&problem, sqlString, problemId); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && problem.UserId != userId && !problem.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO user_favorite_problem (user_id, problem_id, created_at) VALUES ($1, $2, $3)`
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
// @Security ApiKeyAuth
func RemoveProblemFromFavorite(c *gin.Context) {
	problemId := c.Param("id")
	sqlString := `DELETE FROM user_favorite_problem WHERE user_id = $1 AND problem_id = $2`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), problemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "移除成功")
}

// AddProblemSetToFavorite godoc
// @Schemes http
// @Description 添加题集到收藏夹（只有管理员和题集创建者能添加私有题集）
// @Param id path int true "题集ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/favorite/{id} [post]
// @Security ApiKeyAuth
func AddProblemSetToFavorite(c *gin.Context) {
	var problemSet model.ProblemSet
	problemSetId := c.Param("id")
	userId := c.GetInt("UserId")
	sqlString := `SELECT * FROM problemSet WHERE id = $1`
	if err := global.Database.Get(&problemSet, sqlString, problemSetId); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && problemSet.UserId != userId && !problemSet.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO user_favorite_problem_set (user_id, "problem_set_id", created_at) VALUES ($1, $2, $3)`
	if _, err := global.Database.Exec(sqlString, userId, problemSetId, time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveProblemSetFromFavorite godoc
// @Schemes http
// @Description 从收藏夹移除题集
// @Param id path int true "题集ID"
// @Success 200 {string} string "移除成功"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/unfavorite/{id} [post]
// @Security ApiKeyAuth
func RemoveProblemSetFromFavorite(c *gin.Context) {
	problemSetId := c.Param("id")
	sqlString := `DELETE FROM user_favorite_problem_set WHERE user_id = $1 AND "problem_set_id" = $2`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), problemSetId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "移除成功")
}

// FavoriteNote godoc
// @Schemes http
// @Description 收藏笔记
// @Param id path int true "笔记ID"
// @Success 200 {string} string "收藏成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/favorite/{id} [post]
// @Security ApiKeyAuth
func FavoriteNote(c *gin.Context) {
	sqlString := `SELECT * FROM note WHERE id = $1`
	var note model.Note
	if err := global.Database.Get(&note, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && note.UserId != c.GetInt("UserId") && !note.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO user_favorite_note (user_id, note_id, created_at) VALUES ($1, $2, $3) ON CONFLICT do update set created_at = $3`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Param("id"), time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "收藏成功")
}

// UnfavoriteNote godoc
// @Schemes http
// @Description 取消收藏笔记
// @Param id path int true "笔记ID"
// @Success 200 {string} string "取消收藏成功"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/unfavorite/{id} [post]
// @Security ApiKeyAuth
func UnfavoriteNote(c *gin.Context) {
	var note model.Note
	sqlString := `SELECT * FROM note WHERE id = $1`
	if err := global.Database.Get(&note, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	sqlString = `DELETE FROM user_favorite_note WHERE user_id = $1 AND note_id = $2`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "取消收藏成功")
}
