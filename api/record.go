package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type WrongRecordResponse struct {
	ProblemId int       `json:"problem_id"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type AllWrongRecordResponse struct {
	TotalCount int                   `json:"total_count"`
	Records    []WrongRecordResponse `json:"records"`
}

// CreateWrongRecord godoc
// @Schemes http
// @Description 创建错题记录（只有管理员和题目创建者能将私有题目加入到错题记录中）（重复创建会增加做错次数）
// @Tags wrongRecord
// @Param id path int true "题目ID"
// @Success 200 {string} string "创建成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /wrong_record/create/{id} [post]
// @Security ApiKeyAuth
func CreateWrongRecord(c *gin.Context) {
	var problem model.ProblemType
	sqlString := `SELECT * FROM problem_type WHERE id = $1`
	if err := global.Database.Get(&problem, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && problem.UserId != c.GetInt("UserId") && !problem.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO user_wrong_record (user_id, problem_id, count, created_at, updated_at) VALUES ($1, $2, 1, $3, $4) ON CONFLICT 
		(user_id, problem_id) DO UPDATE SET count = user_wrong_record.count + 1, updated_at = $3`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Param("id"), time.Now().Local(), time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

// DeleteWrongRecord godoc
// @Schemes http
// @Description 删除错题记录
// @Tags wrongRecord
// @Param id path int true "题目ID"
// @Success 200 {string} string "删除成功"
// @Failure default {string} string "服务器错误"
// @Router /wrong_record/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteWrongRecord(c *gin.Context) {
	sqlString := `DELETE FROM user_wrong_record WHERE user_id = $1 AND problem_id = $2`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

// GetWrongRecord godoc
// @Schemes http
// @Description 获取错题记录
// @Tags wrongRecord
// @Param id path int true "题目ID"
// @Success 200 {object} AllWrongRecordResponse "获取成功"
// @Failure default {string} string "服务器错误"
// @Router /wrong_record/get/{id} [get]
// @Security ApiKeyAuth
func GetWrongRecord(c *gin.Context) {
	var wrongRecord []model.WrongRecord
	sqlString := `SELECT * FROM user_wrong_record WHERE problem_id = $1 ORDER BY updated_at DESC`
	if err := global.Database.Select(&wrongRecord, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var records []WrongRecordResponse
	for _, record := range wrongRecord {
		records = append(records, WrongRecordResponse{
			ProblemId: record.ProblemId,
			Count:     record.Count,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, AllWrongRecordResponse{
		TotalCount: len(records),
		Records:    records,
	})
}
