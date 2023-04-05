package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type ProblemSetResponse struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	ProblemCount int       `json:"problem_count" db:"problem_count"`
	IsFavorite   bool      `json:"is_favorite" db:"is_favorite"`
}
type ProblemSetRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AllProblemSetResponse struct {
	TotalCount int                  `json:"total_count"`
	ProblemSet []ProblemSetResponse `json:"problemset"`
}

// GetProblemsets godoc
// @Schemes http
// @Description 获取当前用户视角下的所有题集
// @Param id query int false "题集ID"
// @Success 200 {object} AllProblemSetResponse "题集列表"
// @Failure default {string} string "服务器错误"
// @Router /problemset/all [get]
// @Security ApiKeyAuth
func GetProblemsets(c *gin.Context) {
	var problemsets []model.ProblemSet
	var sqlString string
	var err error
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString = `SELECT id, name, description, created_at, updated_at, user_id, is_public FROM problemset WHERE is_public = true`
		if c.Query("id") != "" {
			sqlString += ` AND id = ` + c.Query("id")
		}
		err = global.Database.Select(&problemsets, sqlString)
	} else if role == global.USER {
		sqlString = `SELECT id, name, description, created_at, updated_at, user_id, is_public FROM problemset WHERE (is_public = true OR user_id = $1)`
		if c.Query("id") != "" {
			sqlString += ` AND id = ` + c.Query("id")
		}
		err = global.Database.Select(&problemsets, sqlString, c.GetInt("UserId"))
	} else {
		sqlString = `SELECT id, name, description, created_at, updated_at, user_id, is_public FROM problemset`
		if c.Query("id") != "" {
			sqlString += ` WHERE id = ` + c.Query("id")
		}
		err = global.Database.Select(&problemsets, sqlString)
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var problemsetResponses []ProblemSetResponse
	for _, problemset := range problemsets {
		var problemCount int
		sqlString = `SELECT COUNT(*) FROM problem_in_problemset WHERE problemset_id = $1`
		if err := global.Database.Get(&problemCount, sqlString, problemset.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}

		sqlString = `SELECT COUNT(*) FROM user_favorite_problemset WHERE problemset_id = $1 AND user_id = $2`
		var isFavorite int
		if err := global.Database.Get(&isFavorite, sqlString, problemset.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}

		problemsetResponses = append(problemsetResponses, ProblemSetResponse{
			ID:           problemset.ID,
			Name:         problemset.Name,
			Description:  problemset.Description,
			CreatedAt:    problemset.CreatedAt,
			UpdatedAt:    problemset.UpdatedAt,
			ProblemCount: problemCount,
			IsFavorite:   isFavorite != 0,
		})
	}
	c.JSON(http.StatusOK, AllProblemSetResponse{
		TotalCount: len(problemsetResponses),
		ProblemSet: problemsetResponses,
	})
}

// CreateProblemset godoc
// @Schemes http
// @Description 创建题集
// @Param problemset body ProblemSetRequest true "题集信息"
// @Param is_public query bool true "是否公开"
// @Success 200 {string} string "创建成功"
// @Failure default {string} string "服务器错误"
// @Router /problemset/create [post]
// @Security ApiKeyAuth
func CreateProblemset(c *gin.Context) {
	var problemset ProblemSetRequest
	if err := c.ShouldBindJSON(&problemset); err != nil {
		c.String(http.StatusBadRequest, "请求错误")
		return
	}
	sqlString := `INSERT INTO problemset (name, description, created_at, updated_at, user_id, is_public) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := global.Database.Exec(sqlString, problemset.Name, problemset.Description,
		time.Now().Local(), time.Now().Local(), c.GetInt("UserId"), c.Query("is_public")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

type ProblemResponse struct {
	ID            int `json:"id"`
	ProblemTypeID int `json:"problem_type_id"`
}

// GetProblemsInProblemset godoc
// @Schemes http
// @Description 获取题集中的所有题目信息（只有管理员和题集创建者能获取所有信息，否则只能获取公开题集的所有公开题目信息）
// @Param id path int true "题集ID"
// @Success 200 {object} []ProblemResponse "题目列表"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /problemset/{id}/all_problem [get]
// @Security ApiKeyAuth
func GetProblemsInProblemset(c *gin.Context) {
	var problems []ProblemResponse
	userId := c.GetInt("UserId")
	role, _ := c.Get("Role")
	problemsetId := c.Param("id")
	sqlString := `SELECT is_public, user_id FROM problemset WHERE id = $1`
	var problemsetUserId int
	var isPublic bool
	if err := global.Database.Get(&struct {
		IsPublic bool `db:"is_public"`
		UserId   int  `db:"user_id"`
	}{IsPublic: isPublic, UserId: problemsetUserId}, sqlString, problemsetId); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	} else if !isPublic && userId != problemsetUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	if userId == problemsetUserId || role == global.ADMIN {
		sqlString = `SELECT problem_id, problem_type_id FROM problem_type JOIN problem_in_problemset ON 
    		problem_type.id = problem_in_problemset.problem_id WHERE problemset_id = $1`
		if err := global.Database.Select(&problems, sqlString, problemsetId); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
	} else {
		sqlString = `SELECT problem_id, problem_type_id FROM problem_type JOIN problem_in_problemset ON 
    		problem_type.id = problem_in_problemset.problem_id WHERE problemset_id = $1 AND is_public = true`
		if err := global.Database.Select(&problems, sqlString, problemsetId); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
	}
	c.JSON(http.StatusOK, problems)
}

// AddProblemToProblemset godoc
// @Schemes http
// @Description 添加题目到题集（只有同时为题集的创建者和题目的创建者可以添加题目）
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"/"题目不存在"
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
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if userId != problemsetUserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT user_id FROM problem_type WHERE id = $1`
	var problemUserId int
	if err := global.Database.Get(&problemUserId, sqlString, problemId); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
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
// @Description 从题集中移除题目（只有管理员或者同时为题集的创建者和题目的创建者可以移除题目）
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "移除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"/"题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problemset/{id}/remove [post]
// @Security ApiKeyAuth
func RemoveProblemFromProblemset(c *gin.Context) {
	userId := c.GetInt("UserId")
	role, _ := c.Get("Role")
	problemsetId := c.Param("id")
	problemId := c.Query("problem_id")
	sqlString := `SELECT user_id FROM problemset WHERE id = $1`
	var problemsetUserId int
	if err := global.Database.Get(&problemsetUserId, sqlString, problemsetId); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if userId != problemsetUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT user_id FROM problem_type WHERE id = $1`
	var problemUserId int
	if err := global.Database.Get(&problemUserId, sqlString, problemId); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if userId != problemUserId && role != global.ADMIN {
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
// @Description 删除题集（只有管理员或者题集的创建者可以删除题集）
// @Param id path int true "题集ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /problemset/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteProblemset(c *gin.Context) {
	userId := c.GetInt("UserId")
	problemsetId := c.Param("id")
	sqlString := `SELECT user_id FROM problemset WHERE id = $1`
	var problemsetUserId int
	if err := global.Database.Get(&problemsetUserId, sqlString, problemsetId); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if role, _ := c.Get("Role"); userId != problemsetUserId && role != global.ADMIN {
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
