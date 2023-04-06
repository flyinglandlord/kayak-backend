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
// @Failure default {string} string "已经添加"
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
	if role, _ := c.Get("Role"); role != global.ADMIN && problem.ID != userId && !problem.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO user_favorite_problem (user_id, problem_id, created_at) VALUES ($1, $2, $3)`
	if _, err := global.Database.Exec(sqlString, userId, problemId, time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "已经添加")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveProblemFromFavorite godoc
// @Schemes http
// @Description 从收藏夹移除题目
// @Param id path int true "题目ID"
// @Success 200 {string} string "移除成功"
// @Failure default {string} string "已经移除"
// @Router /problem/unfavorite/{id} [post]
// @Security ApiKeyAuth
func RemoveProblemFromFavorite(c *gin.Context) {
	problemId := c.Param("id")
	userId := c.GetInt("UserId")
	sqlString := `DELETE FROM user_favorite_problem WHERE user_id = $1 AND problem_id = $2`
	if _, err := global.Database.Exec(sqlString, userId, problemId); err != nil {
		c.String(http.StatusInternalServerError, "已经移除")
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
// @Failure default {string} string "已经添加"
// @Router /problemSet/favorite/{id} [post]
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
	sqlString = `INSERT INTO user_favorite_problemSet (user_id, "problemSet_id", created_at) VALUES ($1, $2, $3)`
	if _, err := global.Database.Exec(sqlString, userId, problemSetId, time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "已经添加")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveProblemSetFromFavorite godoc
// @Schemes http
// @Description 从收藏夹移除题集
// @Param id path int true "题集ID"
// @Success 200 {string} string "移除成功"
// @Failure default {string} string "已经移除
// @Router /problemSet/unfavorite/{id} [post]
// @Security ApiKeyAuth
func RemoveProblemSetFromFavorite(c *gin.Context) {
	problemSetId := c.Param("id")
	userId := c.GetInt("UserId")
	sqlString := `DELETE FROM user_favorite_problemSet WHERE user_id = $1 AND "problemSet_id" = $2`
	if _, err := global.Database.Exec(sqlString, userId, problemSetId); err != nil {
		c.String(http.StatusInternalServerError, "已经移除")
		return
	}
	c.String(http.StatusOK, "移除成功")
}
