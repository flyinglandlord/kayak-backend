package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type WrongRecordResponse struct {
	ProblemID int `json:"problem_id"`
	Count     int `json:"count"`
	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
}

// CreateWrongRecord godoc
// @Schemes http
// @Description 创建错题记录（只有管理员和题目创建者能将私有题目加入到错题记录中）（重复创建会增加做错次数）
// @Success 200 {string} string "创建成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /wrong_record/create/{id} [post]
// @Security ApiKeyAuth
func CreateWrongRecord(c *gin.Context) {
	var problem model.ProblemType
	problemID := c.Param("id")
	sqlString := `SELECT * FROM problem_type WHERE id = $1`
	if err := global.Database.Get(&problem, sqlString, problemID); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && problem.UserId != c.GetInt("UserId") && !problem.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO user_wrong_record (user_id, problem_id, count, created_at, updated_at) VALUES ($1, $2, 1, $3, $4) ON CONFLICT 
		(user_id, problem_id) DO UPDATE SET count = user_wrong_record.count + 1, updated_at = $3`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), problemID, time.Now().Local(), time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

// DeleteWrongRecord godoc
// @Schemes http
// @Description 删除错题记录
// @Success 200 {string} string "删除成功"
// @Failure default {string} string "服务器错误"
// @Router /wrong_record/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteWrongRecord(c *gin.Context) {
	problemID := c.Param("id")
	sqlString := `DELETE FROM user_wrong_record WHERE user_id = $1 AND problem_id = $2`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), problemID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}
