package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"strconv"
	"time"
)

type ProblemSetFilter struct {
	ID int `json:"id" db:"id"`
}
type ProblemSetResponse struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	ProblemCount int       `json:"problem_count" db:"problem_count"`
	IsFavorite   bool      `json:"is_favorite" db:"is_favorite"`
}
type AllProblemSetResponse struct {
	TotalCount int                  `json:"total_count"`
	ProblemSet []ProblemSetResponse `json:"problemSet"`
}
type ProblemSetCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

// GetProblemSets godoc
// @Schemes http
// @Description 获取符合题集过滤器要求的当前用户视角下的所有题集（通过表单传输要求）
// @Param filter query ProblemSetFilter false "题集过滤器"
// @Success 200 {object} AllProblemSetResponse "题集列表"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /problemSet/all [get]
// @Security ApiKeyAuth
func GetProblemSets(c *gin.Context) {
	var problemSets []model.ProblemSet
	var filter ProblemSetFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT * FROM problemSet`
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString += ` WHERE is_public = true`
	} else if role == global.USER {
		sqlString += ` WHERE (is_public = true OR user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
	} else {
		sqlString += ` WHERE 1 = 1`
	}
	if filter.ID != 0 {
		sqlString += ` AND id = ` + strconv.Itoa(filter.ID)
	}
	if err := global.Database.Select(&problemSets, sqlString); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var problemSetResponses []ProblemSetResponse
	for _, problemSet := range problemSets {
		var problemCount int
		sqlString = `SELECT COUNT(*) FROM problem_in_problemSet WHERE "problemSet_id" = $1`
		if err := global.Database.Get(&problemCount, sqlString, problemSet.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		sqlString = `SELECT COUNT(*) FROM user_favorite_problemSet WHERE "problemSet_id" = $1 AND user_id = $2`
		var favoriteCount int
		if err := global.Database.Get(&favoriteCount, sqlString, problemSet.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		problemSetResponses = append(problemSetResponses, ProblemSetResponse{
			ID:           problemSet.ID,
			Name:         problemSet.Name,
			Description:  problemSet.Description,
			CreatedAt:    problemSet.CreatedAt,
			UpdatedAt:    problemSet.UpdatedAt,
			ProblemCount: problemCount,
			IsFavorite:   favoriteCount > 0,
		})
	}
	c.JSON(http.StatusOK, AllProblemSetResponse{
		TotalCount: len(problemSetResponses),
		ProblemSet: problemSetResponses,
	})
}

// CreateProblemSet godoc
// @Schemes http
// @Description 创建题集
// @Param problemSet body ProblemSetCreateRequest true "题集信息"
// @Success 200 {string} string "创建成功"
// @Failure 400 {string} string "请求格式错误"
// @Failure default {string} string "服务器错误"
// @Router /problemSet/create [post]
// @Security ApiKeyAuth
func CreateProblemSet(c *gin.Context) {
	var problemSet ProblemSetCreateRequest
	if err := c.ShouldBindJSON(&problemSet); err != nil {
		c.String(http.StatusBadRequest, "请求格式错误")
		return
	}
	sqlString := `INSERT INTO problemSet (name, description, created_at, updated_at, user_id, is_public) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := global.Database.Exec(sqlString, problemSet.Name, problemSet.Description, time.Now().Local(),
		time.Now().Local(), c.GetInt("UserId"), problemSet.IsPublic); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

type ProblemResponse struct {
	ID            int    `json:"id"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	UserId        int    `json:"user_id"`
	IsPublic      bool   `json:"is_public"`
	ProblemTypeID int    `json:"problem_type_id"`
}

// GetProblemsInProblemSet godoc
// @Schemes http
// @Description 获取题集中的所有题目信息（只有管理员和题集创建者能获取所有信息，否则只能获取公开题集的所有公开题目信息）
// @Param id path int true "题集ID"
// @Success 200 {object} []ProblemResponse "题目列表"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /problemSet/{id}/all_problem [get]
// @Security ApiKeyAuth
func GetProblemsInProblemSet(c *gin.Context) {
	var problemSet model.ProblemSet
	var problems []model.ProblemType
	userId := c.GetInt("UserId")
	role, _ := c.Get("Role")
	problemSetId := c.Param("id")
	sqlString := `SELECT * FROM problemSet WHERE id = $1`
	if err := global.Database.Get(&problemSet, sqlString, problemSetId); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if !problemSet.IsPublic && problemSet.UserId != userId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	if userId == problemSet.UserId || role == global.ADMIN {
		sqlString = `SELECT * FROM problem_type pt JOIN problem_in_problemSet pip ON 
    		pt.id = pip.problem_id WHERE "problemSet_id" = $1`
	} else {
		sqlString = `SELECT * FROM problem_type pt JOIN problem_in_problemSet pip ON 
    		pt.id = pip.problem_id WHERE "problemSet_id" = $1 AND pt.is_public = true`
	}
	if err := global.Database.Select(&problems, sqlString, problemSetId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var problemResponses []ProblemResponse
	for _, problem := range problems {
		problemResponses = append(problemResponses, ProblemResponse{
			ID:            problem.ID,
			Description:   problem.Description,
			CreatedAt:     problem.CreatedAt.Format("2006-01-02"),
			UpdatedAt:     problem.UpdatedAt.Format("2006-01-02"),
			UserId:        problem.UserId,
			IsPublic:      problem.IsPublic,
			ProblemTypeID: problem.ProblemTypeId,
		})
	}
	c.JSON(http.StatusOK, problems)
}

// AddProblemToProblemSet godoc
// @Schemes http
// @Description 添加题目到题集（只有同时为题集的创建者和题目的创建者可以添加题目）
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"/"题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problemSet/{id}/add [post]
// @Security ApiKeyAuth
func AddProblemToProblemSet(c *gin.Context) {
	var problemSet model.ProblemSet
	userId := c.GetInt("UserId")
	problemSetId := c.Param("id")
	problemId := c.Query("problem_id")
	sqlString := `SELECT * FROM problemSet WHERE id = $1`
	if err := global.Database.Get(&problemSet, sqlString, problemSetId); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if userId != problemSet.UserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	var problem model.ProblemType
	sqlString = `SELECT * FROM problem_type WHERE id = $1`
	if err := global.Database.Get(&problem, sqlString, problemId); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if userId != problem.UserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO problem_in_problemSet ("problemSet_id", problem_id) VALUES ($1, $2)`
	if _, err := global.Database.Exec(sqlString, problemSetId, problemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveProblemFromProblemSet godoc
// @Schemes http
// @Description 从题集中移除题目（只有管理员或者同时为题集的创建者和题目的创建者可以移除题目）
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "移除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"/"题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problemSet/{id}/remove [post]
// @Security ApiKeyAuth
func RemoveProblemFromProblemSet(c *gin.Context) {
	var problemSet model.ProblemSet
	userId := c.GetInt("UserId")
	role, _ := c.Get("Role")
	problemSetId := c.Param("id")
	problemId := c.Query("problem_id")
	sqlString := `SELECT * FROM problemSet WHERE id = $1`
	if err := global.Database.Get(&problemSet, sqlString, problemSetId); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if userId != problemSet.UserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	var problem model.ProblemType
	sqlString = `SELECT * FROM problem_type WHERE id = $1`
	if err := global.Database.Get(&problem, sqlString, problemId); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if userId != problem.UserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM problem_in_problemSet WHERE "problemSet_id" = $1 AND problem_id = $2`
	if _, err := global.Database.Exec(sqlString, problemSetId, problemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "移除成功")
}

// DeleteProblemSet godoc
// @Schemes http
// @Description 删除题集（只有管理员或者题集的创建者可以删除题集）
// @Param id path int true "题集ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /problemSet/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteProblemSet(c *gin.Context) {
	var problemSet model.ProblemSet
	userId := c.GetInt("UserId")
	problemSetId := c.Param("id")
	sqlString := `SELECT * FROM problemSet WHERE id = $1`
	if err := global.Database.Get(&problemSet, sqlString, problemSetId); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if role, _ := c.Get("Role"); userId != problemSet.UserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM problemSet WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, problemSetId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}
