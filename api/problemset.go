package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"net/http"
	"time"
)

type ProblemsetResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type ProblemsetRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetUserProblemsets godoc
// @Schemes http
// @Description 获取当前登录用户的所有题集
// @Success 200 {object} []ProblemsetResponse "题集列表"
// @Failure default {string} string "服务器错误"
// @Router /user/problemset [get]
// @Security ApiKeyAuth
func GetUserProblemsets(c *gin.Context) {
	var problemsets []ProblemsetResponse
	sqlString := `SELECT id, name, description FROM problemset WHERE user_id = $1`
	if err := global.Database.Select(&problemsets, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, problemsets)
}

// GetProblemsets godoc
// @Schemes http
// @Description 获取当前用户视角下的所有题集
// @Success 200 {object} []ProblemsetResponse "题集列表"
// @Failure default {string} string "服务器错误"
// @Router /problemset/all [get]
func GetProblemsets(c *gin.Context) {
	var problemsets []ProblemsetResponse
	var sqlString string
	var err error
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString = `SELECT id, name, description FROM problemset WHERE is_public = true`
		err = global.Database.Select(&problemsets, sqlString)
	} else if role == global.USER {
		sqlString = `SELECT id, name, description FROM problemset WHERE is_public = true OR user_id = $1`
		err = global.Database.Select(&problemsets, sqlString, c.GetInt("UserId"))
	} else {
		sqlString = `SELECT id, name, description FROM problemset`
		err = global.Database.Select(&problemsets, sqlString)
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, problemsets)
}

// CreateProblemset godoc
// @Schemes http
// @Description 创建题集
// @Param problemset body ProblemsetRequest true "题集信息"
// @Param is_public query bool true "是否公开"
// @Success 200 {string} string "创建成功"
// @Failure default {string} string "服务器错误"
// @Router /problemset/create [post]
// @Security ApiKeyAuth
func CreateProblemset(c *gin.Context) {
	var problemset ProblemsetRequest
	if err := c.ShouldBindJSON(&problemset); err != nil {
		c.String(http.StatusBadRequest, "请求错误")
		return
	}
	sqlString := `INSERT INTO problemset (name, description, created_at, updated_at, user_id, is_public) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := global.Database.Exec(sqlString, problemset.Name, problemset.Description,
		time.Now(), time.Now(), c.GetInt("UserId"), c.Query("is_public")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

// AddProblemToProblemset godoc
// @Schemes http
// @Description 添加题目到题集
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /problemset/{id}/add [post]
// @Security ApiKeyAuth
func AddProblemToProblemset(c *gin.Context) {
	userId := c.GetInt("UserId")
	problemsetId := c.Param("id")
	problemId := c.Query("problem_id")
	sqlString := `SELECT user_id FROM problemset WHERE id = $1`
	var problemsetUserId int
	if err := global.Database.Get(&problemsetUserId, sqlString, problemsetId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if userId != problemsetUserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT user_id FROM problem_type WHERE id = $1`
	var problemUserId int
	if err := global.Database.Get(&problemUserId, sqlString, problemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if userId != problemUserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO problem_in_problemset (problemset_id, problem_id) VALUES ($1, $2)`
	if _, err := global.Database.Exec(sqlString, problemsetId, problemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveProblemFromProblemset godoc
// @Schemes http
// @Description 从题集中移除题目
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "移除成功"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /problemset/{id}/remove [post]
// @Security ApiKeyAuth
func RemoveProblemFromProblemset(c *gin.Context) {
	userId := c.GetInt("UserId")
	problemsetId := c.Param("id")
	problemId := c.Query("problem_id")
	sqlString := `SELECT user_id FROM problemset WHERE id = $1`
	var problemsetUserId int
	if err := global.Database.Get(&problemsetUserId, sqlString, problemsetId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if userId != problemsetUserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM problem_in_problemset WHERE problemset_id = $1 AND problem_id = $2`
	if _, err := global.Database.Exec(sqlString, problemsetId, problemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "移除成功")
}

// DeleteProblemset godoc
// @Schemes http
// @Description 删除题集
// @Param id path int true "题集ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /problemset/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteProblemset(c *gin.Context) {
	userId := c.GetInt("UserId")
	problemsetId := c.Param("id")
	sqlString := `SELECT user_id FROM problemset WHERE id = $1`
	var problemsetUserId int
	if err := global.Database.Get(&problemsetUserId, sqlString, problemsetId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if userId != problemsetUserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM problemset WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, problemsetId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}
