package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type ProblemSetFilter struct {
	ID         *int  `json:"id" form:"id"`
	UserId     *int  `json:"user_id" form:"user_id"`
	IsFavorite *bool `json:"is_favorite" form:"is_favorite"`
	Contain    *int  `json:"contain" form:"contain"`
}
type ProblemSetResponse struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	ProblemCount  int       `json:"problem_count" db:"problem_count"`
	IsFavorite    bool      `json:"is_favorite" db:"is_favorite"`
	FavoriteCount int       `json:"favorite_count" db:"favorite_count"`
	UserId        int       `json:"user_id" db:"user_id"`
	IsPublic      bool      `json:"is_public" db:"is_public"`
}
type ProblemSetRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}
type AllProblemSetResponse struct {
	TotalCount int                  `json:"total_count"`
	ProblemSet []ProblemSetResponse `json:"problem_set"`
}

// GetProblemSets godoc
// @Schemes http
// @Description 获取符合filter要求的当前用户视角下的所有题集
// @Tags problemSet
// @Param filter query ProblemSetFilter false "筛选条件"
// @Success 200 {object} AllProblemSetResponse "题集列表"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/all [get]
// @Security ApiKeyAuth
func GetProblemSets(c *gin.Context) {
	sqlString := `SELECT * FROM problem_set`
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString += ` WHERE is_public = true`
	} else if role == global.USER {
		sqlString += ` WHERE (is_public = true OR user_id = ` + fmt.Sprintf("%d", c.GetInt("UserId")) + `)`
	} else {
		sqlString += ` WHERE 1 = 1`
	}
	var filter ProblemSetFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	if filter.ID != nil {
		sqlString += fmt.Sprintf(` AND id = %d`, *filter.ID)
	}
	if filter.UserId != nil {
		sqlString += fmt.Sprintf(` AND user_id = %d`, *filter.UserId)
	}
	if filter.IsFavorite != nil {
		if *filter.IsFavorite {
			sqlString += ` AND id IN (SELECT problem_set_id FROM user_favorite_problem_set WHERE user_id = ` + fmt.Sprintf("%d", c.GetInt("UserId")) + `)`
		} else {
			sqlString += ` AND id NOT IN (SELECT problem_set_id FROM user_favorite_problem_set WHERE user_id = ` + fmt.Sprintf("%d", c.GetInt("UserId")) + `)`
		}
	}
	if filter.Contain != nil {
		sqlString += ` AND id IN (SELECT problem_set_id FROM problem_in_problemset WHERE problem_id = ` + fmt.Sprintf("%d", *filter.Contain) + `)`
	}
	var problemSets []model.ProblemSet
	if err := global.Database.Select(&problemSets, sqlString); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var problemSetResponses []ProblemSetResponse
	for _, problemSet := range problemSets {
		var problemCount int
		sqlString = `SELECT COUNT(*) FROM problem_in_problem_set WHERE problem_set_id = $1`
		if err := global.Database.Get(&problemCount, sqlString, problemSet.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem_set WHERE problem_set_id = $1 AND user_id = $2`
		var isFavorite int
		if err := global.Database.Get(&isFavorite, sqlString, problemSet.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		var favoriteCount int
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem_set WHERE problem_set_id = $1`
		if err := global.Database.Get(&favoriteCount, sqlString, problemSet.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		problemSetResponses = append(problemSetResponses, ProblemSetResponse{
			ID:            problemSet.ID,
			Name:          problemSet.Name,
			Description:   problemSet.Description,
			CreatedAt:     problemSet.CreatedAt,
			UpdatedAt:     problemSet.UpdatedAt,
			ProblemCount:  problemCount,
			IsFavorite:    isFavorite > 0,
			FavoriteCount: favoriteCount,
			UserId:        problemSet.UserId,
			IsPublic:      problemSet.IsPublic,
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
// @Tags problemSet
// @Param problem_set body ProblemSetRequest true "题集信息"
// @Success 200 {object} ProblemSetResponse "题集信息"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/create [post]
// @Security ApiKeyAuth
func CreateProblemSet(c *gin.Context) {
	var request ProblemSetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求错误")
		return
	}
	tx := global.Database.MustBegin()
	sqlString := `INSERT INTO problem_set (name, description, created_at, updated_at, user_id, is_public) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var problemSetId int
	if err := global.Database.Get(&problemSetId, sqlString, request.Name, request.Description,
		time.Now().Local(), time.Now().Local(), c.GetInt("UserId"), request.IsPublic); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `SELECT * FROM problem_set WHERE id = $1`
	var problemSet model.ProblemSet
	if err := global.Database.Get(&problemSet, sqlString, problemSetId); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if err := tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, ProblemSetResponse{
		ID:            problemSet.ID,
		Name:          problemSet.Name,
		Description:   problemSet.Description,
		CreatedAt:     problemSet.CreatedAt,
		UpdatedAt:     problemSet.UpdatedAt,
		ProblemCount:  0,
		IsFavorite:    false,
		FavoriteCount: 0,
		UserId:        problemSet.UserId,
		IsPublic:      problemSet.IsPublic,
	})
}

type ProblemInProblemSetFilter struct {
	IsFavorite    *bool `json:"is_favorite"`
	ProblemTypeId *int  `json:"problem_type_id"`
}
type ProblemResponse struct {
	ID            int       `json:"id"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserId        int       `json:"user_id"`
	IsPublic      bool      `json:"is_public"`
	IsFavorite    bool      `json:"is_favorite"`
	FavoriteCount int       `json:"favorite_count"`
	ProblemTypeId int       `json:"problem_type_id"`
}
type AllProblemResponse struct {
	TotalCount int               `json:"total_count"`
	Problems   []ProblemResponse `json:"problems"`
}

// GetProblemsInProblemSet godoc
// @Schemes http
// @Description 根据filter获取题集中的所有题目信息
// @Tags problemSet
// @Param id path int true "题集ID"
// @Param filter query ProblemInProblemSetFilter false "筛选条件"
// @Success 200 {object} AllProblemResponse "题目列表"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/all_problem/{id} [get]
// @Security ApiKeyAuth
func GetProblemsInProblemSet(c *gin.Context) {
	sqlString := `SELECT * FROM problem_set WHERE id = $1`
	var problemSet model.ProblemSet
	if err := global.Database.Get(&problemSet, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && problemSet.UserId != c.GetInt("UserId") && !problemSet.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	var filter ProblemInProblemSetFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString = `SELECT * FROM problem_type` + fmt.Sprintf(" WHERE id IN (SELECT problem_id FROM problem_in_problem_set WHERE problem_set_id = %d)", problemSet.ID)
	if filter.IsFavorite != nil {
		sqlString += fmt.Sprintf(" AND id IN (SELECT problem_id FROM user_favorite_problem WHERE user_id = %d)", c.GetInt("UserId"))
	}
	if filter.ProblemTypeId != nil {
		sqlString += fmt.Sprintf(" AND problem_type_id = %d", *filter.ProblemTypeId)
	}
	var problems []model.ProblemType
	if err := global.Database.Select(&problems, sqlString); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var problemResponses []ProblemResponse
	for _, problem := range problems {
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem WHERE problem_id = $1 AND user_id = $2`
		var isFavorite int
		if err := global.Database.Get(&isFavorite, sqlString, problem.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem WHERE problem_id = $1`
		var favoriteCount int
		if err := global.Database.Get(&favoriteCount, sqlString, problem.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		problemResponses = append(problemResponses, ProblemResponse{
			ID:            problem.ID,
			Description:   problem.Description,
			CreatedAt:     problem.CreatedAt,
			UpdatedAt:     problem.UpdatedAt,
			UserId:        problem.UserId,
			IsPublic:      problem.IsPublic,
			IsFavorite:    isFavorite > 0,
			FavoriteCount: favoriteCount,
			ProblemTypeId: problem.ProblemTypeId,
		})
	}
	c.JSON(http.StatusOK, AllProblemResponse{
		TotalCount: len(problemResponses),
		Problems:   problemResponses,
	})
}

// AddProblemToProblemSet godoc
// @Schemes http
// @Description 添加题目到题集（只有同时为题集的创建者和题目的创建者可以添加题目）
// @Tags problemSet
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"/"题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/add/{id} [post]
// @Security ApiKeyAuth
func AddProblemToProblemSet(c *gin.Context) {
	sqlString := `SELECT user_id FROM problem_set WHERE id = $1`
	var problemSetUserId int
	if err := global.Database.Get(&problemSetUserId, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if c.GetInt("UserId") != problemSetUserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT user_id FROM problem_type WHERE id = $1`
	var problemUserId int
	if err := global.Database.Get(&problemUserId, sqlString, c.Query("problem_id")); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if c.GetInt("UserId") != problemUserId {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO problem_in_problem_set (problem_set_id, problem_id) VALUES ($1, $2)`
	if _, err := global.Database.Exec(sqlString, c.Param("id"), c.Query("problem_id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveProblemFromProblemSet godoc
// @Schemes http
// @Description 从题集中移除题目（只有管理员或者同时为题集的创建者和题目的创建者可以移除题目）
// @Tags problemSet
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "移除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"/"题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/remove/{id} [delete]
// @Security ApiKeyAuth
func RemoveProblemFromProblemSet(c *gin.Context) {
	role, _ := c.Get("Role")
	sqlString := `SELECT user_id FROM problem_set WHERE id = $1`
	var problemSetUserId int
	if err := global.Database.Get(&problemSetUserId, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if c.GetInt("UserId") != problemSetUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT user_id FROM problem_type WHERE id = $1`
	var problemUserId int
	if err := global.Database.Get(&problemUserId, sqlString, c.Query("problem_id")); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if c.GetInt("UserId") != problemUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM problem_in_problem_set WHERE problem_set_id = $1 AND problem_id = $2`
	if _, err := global.Database.Exec(sqlString, c.Param("id"), c.Query("problem_id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "移除成功")
}

// DeleteProblemSet godoc
// @Schemes http
// @Description 删除题集（只有管理员或者题集的创建者可以删除题集）
// @Tags problemSet
// @Param id path int true "题集ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteProblemSet(c *gin.Context) {
	sqlString := `SELECT user_id FROM problem_set WHERE id = $1`
	var problemSetUserId int
	if err := global.Database.Get(&problemSetUserId, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	if role, _ := c.Get("Role"); c.GetInt("UserId") != problemSetUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM problem_set WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}
