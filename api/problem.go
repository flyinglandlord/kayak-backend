package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"net/http"
	"time"
)

const (
	ChoiceProblemType = iota
	BlankProblemType
)

func DeleteProblem(c *gin.Context) {
	userId := c.GetInt("UserId")
	choiceProblemId := c.Param("id")
	sqlString := `SELECT user_id FROM problem_type WHERE id = $1`
	var problemUserId int
	if err := global.Database.Get(&problemUserId, sqlString, choiceProblemId); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if role, _ := c.Get("Role"); userId != problemUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM problem_type WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, choiceProblemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

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
	IsCorrect   bool   `json:"-" db:"is_correct"`
}

// GetChoiceProblems godoc
// @Schemes http
// @Description 获取当前用户视角下的所有选择题
// @Success 200 {object} []ChoiceProblemResponse "选择题列表"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice/all [get]
// @Security ApiKeyAuth
func GetChoiceProblems(c *gin.Context) {
	var choiceProblems []ChoiceProblemResponse
	var sqlString string
	var err error
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString = `SELECT id, description FROM problem_type WHERE problem_type_id = $1 AND is_public = true`
		err = global.Database.Select(&choiceProblems, sqlString, ChoiceProblemType)
	} else if role == global.USER {
		sqlString = `SELECT id, description FROM problem_type WHERE problem_type_id = $1 AND (is_public = true OR user_id = $2)`
		err = global.Database.Select(&choiceProblems, sqlString, ChoiceProblemType, c.GetInt("UserId"))
	} else {
		sqlString = `SELECT id, description FROM problem_type WHERE problem_type_id = $1`
		err = global.Database.Select(&choiceProblems, sqlString, ChoiceProblemType)
	}
	if err != nil {
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

// GetChoiceProblem godoc
// @Schemes http
// @Description 获取单个选择题信息（只有管理员和题目创建者可以获取私有题目）
// @Param id path int true "选择题ID"
// @Success 200 {object} ChoiceProblemResponse "选择题信息"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "选择题不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice/{id} [get]
func GetChoiceProblem(c *gin.Context) {
	var choiceProblem ChoiceProblemResponse
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString := `SELECT is_public FROM problem_type WHERE problem_type_id=$1 AND id = $2`
		var isPublic bool
		if err := global.Database.Get(&isPublic, sqlString, ChoiceProblemType, c.Param("id")); err != nil {
			c.String(http.StatusNotFound, "选择题不存在")
			return
		}
		if !isPublic {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	} else if role == global.USER {
		sqlString := `SELECT is_public, user_id FROM problem_type WHERE problem_type_id=$1 AND id = $2`
		var isPublic bool
		var userId int
		if err := global.Database.Get(&struct {
			IsPublic bool `db:"is_public"`
			UserId   int  `db:"user_id"`
		}{IsPublic: isPublic, UserId: userId}, sqlString, ChoiceProblemType, c.Param("id")); err != nil {
			c.String(http.StatusNotFound, "选择题不存在")
			return
		}
		if !isPublic && userId != c.GetInt("UserId") {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	}
	sqlString := `SELECT id, description FROM problem_type WHERE problem_type_id=$1 AND id = $2`
	if err := global.Database.Get(&choiceProblem, sqlString, ChoiceProblemType, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "选择题不存在")
		return
	}
	sqlString = `SELECT choice, description, is_correct FROM problem_choice WHERE id = $1`
	if err := global.Database.Select(&choiceProblem.Choices, sqlString, choiceProblem.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, choiceProblem)
}

// CreateChoiceProblem godoc
// @Schemes http
// @Description 创建选择题
// @Param problem body ChoiceProblemRequest true "选择题信息"
// @Param is_public query bool true "是否公开"
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
		ChoiceProblemType, c.Query("is_public"), time.Now().Local(), time.Now().Local()); err != nil {
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
// @Description 更新选择题（description为空或不传则维持原description,choices只需传需要更新的）(只有管理员和题目创建者可以更新题目)
// @Param problem body ChoiceProblemResponse true "选择题信息"
// @Param is_public query bool true "是否公开"
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
	if role, _ := c.Get("Role"); userId != problemUserId && role != global.ADMIN {
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
		c.Query("is_public"), time.Now().Local(), choiceProblem.ID); err != nil {
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
// @Description 删除选择题（只有管理员和题目创建者可以删除题目）
// @Param id path int true "选择题ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteChoiceProblem(c *gin.Context) {
	DeleteProblem(c)
}

type BlankProblemResponse struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}
type BlankProblemRequest struct {
	Description string `json:"description"`
}

// GetBlankProblems godoc
// @Schemes http
// @Description 获取当前用户视角下的所有填空题
// @Success 200 {object} BlankProblemResponse "填空题信息"
// @Failure default {string} string "服务器错误"
// @Router /problem/blank/all [get]
// @Security ApiKeyAuth
func GetBlankProblems(c *gin.Context) {
	var blankProblems []BlankProblemResponse
	var sqlString string
	var err error
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString = `SELECT id, description FROM problem_type WHERE problem_type_id = $1 AND is_public = true`
		err = global.Database.Select(&blankProblems, sqlString, BlankProblemType)
	} else if role == global.USER {
		sqlString = `SELECT id, description FROM problem_type WHERE problem_type_id = $1 AND (is_public = true OR user_id = $2)`
		err = global.Database.Select(&blankProblems, sqlString, BlankProblemType, c.GetInt("UserId"))
	} else {
		sqlString = `SELECT id, description FROM problem_type WHERE problem_type_id = $1`
		err = global.Database.Select(&blankProblems, sqlString, BlankProblemType)
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, blankProblems)
}

// GetBlankProblem godoc
// @Schemes http
// @Description 获取单个填空题信息（只有管理员和题目创建者可以查看私有题目）
// @Param id path int true "填空题ID"
// @Success 200 {object} BlankProblemResponse "填空题信息"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "选择题不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/blank/:id [get]
func GetBlankProblem(c *gin.Context) {
	var blankProblem BlankProblemResponse
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString := `SELECT is_public FROM problem_type WHERE problem_type_id=$1 AND id = $2`
		var isPublic bool
		if err := global.Database.Get(&isPublic, sqlString, BlankProblemType, c.Param("id")); err != nil {
			c.String(http.StatusNotFound, "填空题不存在")
			return
		}
		if !isPublic {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	} else if role == global.USER {
		sqlString := `SELECT is_public, user_id FROM problem_type WHERE problem_type_id=$1 AND id = $2`
		var isPublic bool
		var userId int
		if err := global.Database.Get(&struct {
			IsPublic bool `db:"is_public"`
			UserId   int  `db:"user_id"`
		}{IsPublic: isPublic, UserId: userId}, sqlString, BlankProblemType, c.Param("id")); err != nil {
			c.String(http.StatusNotFound, "填空题不存在")
			return
		}
		if !isPublic && userId != c.GetInt("UserId") {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	}
	sqlString := `SELECT id, description FROM problem_type WHERE problem_type_id=$1 AND id = $2`
	if err := global.Database.Get(&blankProblem, sqlString, BlankProblemType, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "填空题不存在")
		return
	}
	c.JSON(http.StatusOK, blankProblem)
}

// CreateBlankProblem godoc
// @Schemes http
// @Description 创建填空题
// @Param problem body BlankProblemRequest true "填空题信息"
// @Param is_public query bool true "是否公开"
// @Success 200 {string} string "创建成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /problem/blank/create [post]
// @Security ApiKeyAuth
func CreateBlankProblem(c *gin.Context) {
	userId := c.GetInt("UserId")
	var blankProblem BlankProblemRequest
	if err := c.ShouldBindJSON(&blankProblem); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `INSERT INTO problem_type (problem_type_id, description, is_public, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := global.Database.Exec(sqlString, BlankProblemType, blankProblem.Description,
		c.Query("is_public"), userId, time.Now().Local(), time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

// UpdateBlankProblem godoc
// @Schemes http
// @Description 更新填空题（只有管理员和题目创建者可以更新题目）
// @Param problem body BlankProblemResponse true "填空题信息"
// @Param is_public query bool true "是否公开"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /problem/blank/update [put]
// @Security ApiKeyAuth
func UpdateBlankProblem(c *gin.Context) {
	userId := c.GetInt("UserId")
	var blankProblem BlankProblemResponse
	if err := c.ShouldBindJSON(&blankProblem); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT user_id FROM problem_type WHERE id = $1`
	var problemUserId int
	if err := global.Database.Get(&problemUserId, sqlString, blankProblem.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if role, _ := c.Get("Role"); userId != problemUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `UPDATE problem_type SET description = $1, is_public = $2, updated_at = $3 WHERE id = $4`
	if _, err := global.Database.Exec(sqlString, blankProblem.Description, c.Query("is_public"),
		time.Now().Local(), blankProblem.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "更新成功")
}

// DeleteBlankProblem godoc
// @Schemes http
// @Description 删除填空题（只有管理员和题目创建者可以删除题目）
// @Param id path int true "填空题ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/blank/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteBlankProblem(c *gin.Context) {
	DeleteProblem(c)
}
