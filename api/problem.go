package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"net/http"
	"time"
)

const (
	ChoiceProblemType = iota
)

type ChoiceProblemResponse struct {
	ID          int      `json:"id"`
	Description string   `json:"description"`
	Choices     []Choice `json:"choices"`
}
type ChoiceProblemRequest struct {
	Description string   `json:"description"`
	Choices     []Choice `json:"choices"`
}
type Choice struct {
	Choice      string `json:"choice"`
	Description string `json:"description"`
	IsCorrect   bool   `json:"is_correct" db:"is_correct"`
}

// GetChoiceProblems godoc
// @Schemes http
// @Description 获取所有的选择题
// @Success 200 {object} []ChoiceProblemResponse "选择题列表"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice [get]
// @Security ApiKeyAuth
func GetChoiceProblems(c *gin.Context) {
	var choiceProblems []ChoiceProblemResponse
	sqlString := `SELECT id, description FROM problem_type WHERE problem_type_id = $1`
	if err := global.Database.Select(&choiceProblems, sqlString, ChoiceProblemType); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	for i := range choiceProblems {
		sqlString = `SELECT choice, description, is_correct FROM problem_choice WHERE id = $1`
		if err := global.Database.Select(&choiceProblems[i].Choices, sqlString, choiceProblems[i].ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
	}
	c.JSON(http.StatusOK, choiceProblems)
}

// CreateChoiceProblem godoc
// @Schemes http
// @Description 创建选择题
// @Param problem body ChoiceProblemRequest true "选择题信息"
// @Param is_public query bool false "是否公开"
// @Success 200 {string} string "创建成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice/create [post]
// @Security ApiKeyAuth
func CreateChoiceProblem(c *gin.Context) {
	choiceProblem := ChoiceProblemResponse{}
	if err := c.ShouldBindJSON(&choiceProblem); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	tx := global.Database.MustBegin()
	userId := c.GetInt("UserId")
	sqlString := `INSERT INTO problem_type (description, user_id, problem_type_id, is_public, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	if err := global.Database.Get(&choiceProblem.ID, sqlString, choiceProblem.Description, userId,
		ChoiceProblemType, c.Query("is_public"), time.Now(), time.Now()); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	for i := range choiceProblem.Choices {
		sqlString = `INSERT INTO problem_choice (id, choice, description, is_correct) VALUES ($1, $2, $3, $4)`
		if _, err := tx.Exec(sqlString, choiceProblem.ID, choiceProblem.Choices[i].Choice,
			choiceProblem.Choices[i].Description, choiceProblem.Choices[i].IsCorrect); err != nil {
			_ = tx.Rollback()
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
	}
	if err := tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

// UpdateChoiceProblem godoc
// @Schemes http
// @Description 更新选择题(description为空或不传则维持原description,choices只需传需要更新的)
// @Param problem body ChoiceProblemResponse true "选择题信息"
// @Param is_public query bool false "是否公开"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice/update [put]
// @Security ApiKeyAuth
func UpdateChoiceProblem(c *gin.Context) {
	choiceProblem := ChoiceProblemResponse{}
	if err := c.ShouldBindJSON(&choiceProblem); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	userId := c.GetInt("UserId")
	sqlString := `SELECT user_id FROM problem_type WHERE id = $1`
	var problemUserId int
	if err := global.Database.Get(&problemUserId, sqlString, choiceProblem.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if userId != problemUserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	tx := global.Database.MustBegin()
	if choiceProblem.Description == "" {
		sqlString := `SELECT description FROM problem_type WHERE id = $1`
		if err := global.Database.Get(&choiceProblem.Description, sqlString, choiceProblem.ID); err != nil {
			_ = tx.Rollback()
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
	}
	sqlString = `UPDATE problem_type SET description = $1, is_public = $2, updated_at = $3 WHERE id = $4`
	if _, err := global.Database.Exec(sqlString, choiceProblem.Description,
		c.Query("is_public"), time.Now(), choiceProblem.ID); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	for i := range choiceProblem.Choices {
		sqlString = `INSERT INTO problem_choice (id, choice, description, is_correct) VALUES ($1, $2, $3, $4) 
			ON CONFLICT (id, choice) DO UPDATE SET description = $3, is_correct = $4`
		if _, err := global.Database.Exec(sqlString, choiceProblem.ID, choiceProblem.Choices[i].Choice,
			choiceProblem.Choices[i].Description, choiceProblem.Choices[i].IsCorrect); err != nil {
			_ = tx.Rollback()
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
	}
	if err := tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "更新成功")
}

// DeleteChoiceProblem godoc
// @Schemes http
// @Description 删除选择题
// @Param id path int true "选择题ID"
// @Success 200 {string} string "删除成功"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteChoiceProblem(c *gin.Context) {
	userId := c.GetInt("UserId")
	choiceProblemId := c.Param("id")
	sqlString := `SELECT user_id FROM problem_type WHERE id = $1`
	var problemUserId int
	if err := global.Database.Get(&problemUserId, sqlString, choiceProblemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if userId != problemUserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	tx := global.Database.MustBegin()
	sqlString = `DELETE FROM problem_choice WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, choiceProblemId); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `DELETE FROM problem_type WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, choiceProblemId); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if err := tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}
