package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"strconv"
	"time"
)

type ProblemSetFilter struct {
	ID         *int  `json:"id" form:"id"`
	UserId     *int  `json:"user_id" form:"user_id"`
	GroupId    *int  `json:"group_id" form:"group_id"`
	IsPublic   *bool `json:"is_public" form:"is_public"`
	IsFavorite *bool `json:"is_favorite" form:"is_favorite"`
	Contain    *int  `json:"contain" form:"contain"`
	AreaId     *int  `json:"area_id" form:"area_id"`
}
type ProblemSetResponse struct {
	ID            int              `json:"id" db:"id"`
	Name          string           `json:"name" db:"name"`
	Description   string           `json:"description" db:"description"`
	CreatedAt     time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" db:"updated_at"`
	ProblemCount  int              `json:"problem_count" db:"problem_count"`
	IsFavorite    bool             `json:"is_favorite" db:"is_favorite"`
	FavoriteCount int              `json:"favorite_count" db:"favorite_count"`
	UserId        int              `json:"user_id" db:"user_id"`
	UserInfo      UserInfoResponse `json:"user_info" db:"user_info"`
	IsPublic      bool             `json:"is_public" db:"is_public"`
	GroupId       int              `json:"group_id" db:"group_id"`
	AreaId        int              `json:"area_id" db:"area_id"`
}
type ProblemSetCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	GroupId     *int   `json:"group_id"`
	AreaId      *int   `json:"area_id"`
}
type AllProblemSetResponse struct {
	TotalCount int                  `json:"total_count"`
	ProblemSet []ProblemSetResponse `json:"problem_set"`
}
type ProblemSetUpdateRequest struct {
	ID          int     `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsPublic    *bool   `json:"is_public"`
	GroupId     *int    `json:"group_id"`
	AreaId      *int    `json:"area_id"`
}

// GetProblemSets godoc
// @Schemes http
// @Description 获取符合filter要求的当前用户视角下的所有题集
// @Tags ProblemSet
// @Param filter query ProblemSetFilter false "筛选条件"
// @Success 200 {object} AllProblemSetResponse "题集列表"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/all [get]
// @Security ApiKeyAuth
func GetProblemSets(c *gin.Context) {
	var filter ProblemSetFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT * FROM problem_set`
	role, _ := c.Get("Role")
	if filter.GroupId == nil {
		if role == global.GUEST {
			sqlString += ` WHERE is_public = true`
		} else if role == global.USER {
			sqlString += ` WHERE (is_public = true OR user_id IN (SELECT user_id FROM group_member WHERE group_id = problem_set.group_id) OR user_id =` + fmt.Sprintf("%d", c.GetInt("UserId")) + `)`
		} else {
			sqlString += ` WHERE 1 = 1`
		}
	} else if *filter.GroupId != 0 {
		if role == global.GUEST {
			sqlString += ` WHERE is_public = true `
		} else if role == global.USER {
			sqlString += fmt.Sprint(` WHERE (is_public = true OR `, c.GetInt("UserId"), ` IN (SELECT user_id FROM group_member WHERE group_id = `, *filter.GroupId, `))`)
		} else {
			sqlString += ` WHERE 1 = 1`
		}
		sqlString += ` AND group_id = ` + fmt.Sprintf("%d", *filter.GroupId)
	} else {
		sqlString += ` WHERE group_id = 0`
	}
	if filter.ID != nil {
		sqlString += fmt.Sprintf(` AND id = %d`, *filter.ID)
	}
	if filter.AreaId != nil {
		sqlString += fmt.Sprintf(` AND area_id = %d`, *filter.AreaId)
	}
	if filter.UserId != nil {
		sqlString += fmt.Sprint(` AND user_id = `, *filter.UserId)
		sqlString += fmt.Sprint(` AND (group_id = 0 OR (`, *filter.UserId, ` IN (SELECT user_id FROM group_member WHERE group_member.group_id = problem_set.group_id)))`)
	}
	if filter.IsFavorite != nil {
		if *filter.IsFavorite {
			sqlString += ` AND id IN (SELECT problem_set_id FROM user_favorite_problem_set WHERE user_id = ` + fmt.Sprintf("%d", c.GetInt("UserId")) + `)`
		} else {
			sqlString += ` AND id NOT IN (SELECT problem_set_id FROM user_favorite_problem_set WHERE user_id = ` + fmt.Sprintf("%d", c.GetInt("UserId")) + `)`
		}
	}
	if filter.IsPublic != nil {
		sqlString += fmt.Sprintf(` AND is_public = %t`, *filter.IsPublic)
	}
	if filter.Contain != nil {
		sqlString += ` AND id IN (SELECT problem_set_id FROM problem_in_problem_set WHERE problem_id = ` + fmt.Sprintf("%d", *filter.Contain) + `)`
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
		user := model.User{}
		sqlString = `SELECT id, avatar_url, nick_name FROM "user" WHERE id = $1`
		if err := global.Database.Get(&user, sqlString, problemSet.UserId); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		userInfo := UserInfoResponse{
			UserId:     user.ID,
			AvatarPath: user.AvatarURL,
			NickName:   user.NickName,
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
			UserInfo:      userInfo,
			IsPublic:      problemSet.IsPublic,
			GroupId:       problemSet.GroupId,
			AreaId:        problemSet.AreaId,
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
// @Tags ProblemSet
// @Param problem_set body ProblemSetCreateRequest true "题集信息"
// @Success 200 {object} ProblemSetResponse "题集信息"
// @Failure 400 {string} string "请求错误"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/create [post]
// @Security ApiKeyAuth
func CreateProblemSet(c *gin.Context) {
	var request ProblemSetCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求错误")
		return
	}
	tx := global.Database.MustBegin()
	var problemSetId int
	if request.GroupId == nil {
		request.GroupId = new(int)
		*request.GroupId = 0
	} else if *request.GroupId != 0 {
		sqlString := `SELECT COUNT(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, request.GroupId, c.GetInt("UserId")); err != nil {
			_ = tx.Rollback()
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		if count == 0 {
			_ = tx.Rollback()
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	}
	if request.AreaId == nil {
		request.AreaId = new(int)
		*request.AreaId = 100
	}
	sqlString := `INSERT INTO problem_set (name, description, created_at, updated_at, user_id, is_public, group_id, area_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	if err := global.Database.Get(&problemSetId, sqlString, request.Name, request.Description, time.Now().Local(),
		time.Now().Local(), c.GetInt("UserId"), request.IsPublic, request.GroupId, request.AreaId); err != nil {
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
		GroupId:       problemSet.GroupId,
		AreaId:        problemSet.AreaId,
	})
}

// UpdateProblemSet godoc
// @Schemes http
// @Description 更新题集(只需传需要更改的)
// @Tags ProblemSet
// @Param problem_set body ProblemSetUpdateRequest true "题集信息"
// @Success 200 {object} string "更新成功"
// @Failure 400 {string} string "请求错误"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/update [put]
// @Security ApiKeyAuth
func UpdateProblemSet(c *gin.Context) {
	var request ProblemSetUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求错误")
		return
	}
	var problemSet model.ProblemSet
	sqlString := `SELECT * FROM problem_set WHERE id = $1`
	if err := global.Database.Get(&problemSet, sqlString, request.ID); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	role, _ := c.Get("Role")
	if problemSet.GroupId == 0 {
		if role != global.ADMIN && problemSet.UserId != c.GetInt("UserId") {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	} else {
		sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, problemSet.GroupId, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		if role != global.ADMIN && count == 0 {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	}
	if request.Name == nil {
		request.Name = &problemSet.Name
	}
	if request.Description == nil {
		request.Description = &problemSet.Description
	}
	if request.IsPublic == nil {
		request.IsPublic = &problemSet.IsPublic
	}
	if request.GroupId == nil {
		request.GroupId = &problemSet.GroupId
	} else if *request.GroupId != 0 {
		sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, request.GroupId, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		if role != global.ADMIN && count == 0 {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	}
	if request.AreaId == nil {
		request.AreaId = &problemSet.AreaId
	}
	sqlString = `UPDATE problem_set SET name = $1, description = $2, updated_at = $3, is_public = $4, group_id = $5, area_id = $6 WHERE id = $7`
	if _, err := global.Database.Exec(sqlString, request.Name, request.Description,
		time.Now().Local(), request.IsPublic, request.GroupId, request.AreaId, request.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "更新成功")
}

type ProblemInProblemSetFilter struct {
	IsFavorite    *bool `json:"is_favorite" form:"is_favorite"`
	ProblemTypeId *int  `json:"problem_type_id" form:"problem_type_id"`
	IsWrong       *bool `json:"is_wrong" form:"is_wrong"`
	Offset        *int  `json:"offset" form:"offset"`
	Limit         *int  `json:"limit" form:"limit"`
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
// @Tags ProblemSet
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
	role, _ := c.Get("Role")
	if problemSet.GroupId == 0 {
		if role != global.ADMIN && problemSet.UserId != c.GetInt("UserId") && !problemSet.IsPublic {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	} else {
		sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, problemSet.GroupId, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		if role != global.ADMIN && count == 0 && !problemSet.IsPublic {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	}
	var filter ProblemInProblemSetFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString = `SELECT * FROM problem_type` + fmt.Sprintf(" WHERE id IN (SELECT problem_id FROM problem_in_problem_set WHERE problem_set_id = %d)", problemSet.ID)
	if filter.IsFavorite != nil {
		if *filter.IsFavorite {
			sqlString += fmt.Sprintf(" AND id IN (SELECT problem_id FROM user_favorite_problem WHERE user_id = %d)", c.GetInt("UserId"))
		} else {
			sqlString += fmt.Sprintf(" AND id NOT IN (SELECT problem_id FROM user_favorite_problem WHERE user_id = %d)", c.GetInt("UserId"))
		}
	}
	if filter.ProblemTypeId != nil {
		sqlString += fmt.Sprintf(" AND problem_type_id = %d", *filter.ProblemTypeId)
	}
	if filter.IsWrong != nil {
		if *filter.IsWrong {
			sqlString += fmt.Sprintf(` AND id IN (SELECT problem_id FROM user_wrong_record WHERE user_id = %d)`, c.GetInt("UserId"))
		} else {
			sqlString += fmt.Sprintf(` AND id NOT IN (SELECT problem_id FROM user_wrong_record WHERE user_id = %d)`, c.GetInt("UserId"))
		}
	}
	if filter.Limit != nil {
		sqlString += ` LIMIT ` + strconv.Itoa(*filter.Limit)
	}
	if filter.Offset != nil {
		sqlString += ` OFFSET ` + strconv.Itoa(*filter.Offset)
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
// @Tags ProblemSet
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"/"题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/add/{id} [post]
// @Security ApiKeyAuth
func AddProblemToProblemSet(c *gin.Context) {
	sqlString := `SELECT * FROM problem_set WHERE id = $1`
	var problemSet model.ProblemSet
	if err := global.Database.Get(&problemSet, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	role, _ := c.Get("Role")
	if problemSet.GroupId == 0 {
		if role != global.ADMIN && problemSet.UserId != c.GetInt("UserId") {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	} else {
		sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, problemSet.GroupId, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		if role != global.ADMIN && count == 0 {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
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

// MigrateProblemToProblemSet godoc
// @Schemes http
// @Description 复制题目到题集
// @Tags ProblemSet
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "复制成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"/"题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/migrate/{id} [post]
// @Security ApiKeyAuth
func MigrateProblemToProblemSet(c *gin.Context) {
	sqlString := `SELECT * FROM problem_set WHERE id = $1`
	var problemSet model.ProblemSet
	if err := global.Database.Get(&problemSet, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	role, _ := c.Get("Role")
	if problemSet.GroupId == 0 {
		if role != global.ADMIN && problemSet.UserId != c.GetInt("UserId") {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	} else {
		sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, problemSet.GroupId, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		if role != global.ADMIN && count == 0 {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	}
	sqlString = `SELECT * FROM problem_type WHERE id = $1`
	var problem model.ProblemType
	if err := global.Database.Get(&problem, sqlString, c.Query("problem_id")); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	var problemSetIds []int
	sqlString = `SELECT problem_set_id FROM problem_in_problem_set WHERE problem_id = $1`
	if err := global.Database.Select(&problemSetIds, sqlString, c.Query("problem_id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var groupId int
	sqlString = `SELECT group_id FROM problem_set WHERE id = $1`
	if err := global.Database.Get(&groupId, sqlString, problemSetIds[0]); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if groupId == 0 {
		if role != global.ADMIN && problem.UserId != c.GetInt("UserId") && !problem.IsPublic {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	} else {
		sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, groupId, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		if role != global.ADMIN && count == 0 && !problem.IsPublic {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	}
	tx := global.Database.MustBegin()
	sqlString = `INSERT INTO problem_type (description, created_at, updated_at, user_id, 
  		is_public, problem_type_id, analysis) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	if err := global.Database.Get(&problem.ID, sqlString, problem.Description, time.Now().Local(), time.Now().Local(),
		c.GetInt("UserId"), problem.IsPublic, problem.ProblemTypeId, problem.Analysis); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `INSERT INTO problem_in_problem_set (problem_set_id, problem_id) VALUES ($1, $2)`
	if _, err := global.Database.Exec(sqlString, c.Param("id"), problem.ID); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if problem.ProblemTypeId == ChoiceProblemType {
		var choices []model.ProblemChoice
		sqlString = `SELECT * FROM problem_choice WHERE id = $1`
		if err := global.Database.Select(&choices, sqlString, c.Query("problem_id")); err != nil {
			_ = tx.Rollback()
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		for _, choice := range choices {
			sqlString = `INSERT INTO problem_choice (id, choice, description, is_correct) VALUES ($1, $2, $3, $4)`
			if _, err := global.Database.Exec(sqlString, problem.ID, choice.Choice, choice.Description, choice.IsCorrect); err != nil {
				_ = tx.Rollback()
				c.String(http.StatusInternalServerError, "服务器错误")
				return
			}
		}
	} else if problem.ProblemTypeId == BlankProblemType {
		var answer model.ProblemAnswer
		sqlString = `SELECT * FROM problem_answer WHERE id = $1`
		if err := global.Database.Get(&answer, sqlString, c.Query("problem_id")); err != nil {
			_ = tx.Rollback()
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		sqlString = `INSERT INTO problem_answer (id, answer) VALUES ($1, $2)`
		if _, err := global.Database.Exec(sqlString, problem.ID, answer.Answer); err != nil {
			_ = tx.Rollback()
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
	} else if problem.ProblemTypeId == JudgeProblemType {
		var judge model.ProblemJudge
		sqlString = `SELECT * FROM problem_judge WHERE id = $1`
		if err := global.Database.Get(&judge, sqlString, c.Query("problem_id")); err != nil {
			_ = tx.Rollback()
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		sqlString = `INSERT INTO problem_judge (id, is_correct) VALUES ($1, $2)`
		if _, err := global.Database.Exec(sqlString, problem.ID, judge.IsCorrect); err != nil {
			_ = tx.Rollback()
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveProblemFromProblemSet godoc
// @Schemes http
// @Description 从题集中移除题目（只有管理员或者同时为题集的创建者和题目的创建者可以移除题目）
// @Tags ProblemSet
// @Param id path int true "题集ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "移除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"/"题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/remove/{id} [delete]
// @Security ApiKeyAuth
func RemoveProblemFromProblemSet(c *gin.Context) {
	sqlString := `SELECT * FROM problem_set WHERE id = $1`
	var problemSet model.ProblemSet
	if err := global.Database.Get(&problemSet, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	role, _ := c.Get("Role")
	if problemSet.GroupId == 0 {
		if role != global.ADMIN && problemSet.UserId != c.GetInt("UserId") {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	} else {
		sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, problemSet.GroupId, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		if role != global.ADMIN && count == 0 {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
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
// @Tags ProblemSet
// @Param id path int true "题集ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteProblemSet(c *gin.Context) {
	sqlString := `SELECT * FROM problem_set WHERE id = $1`
	var problemSet model.ProblemSet
	if err := global.Database.Get(&problemSet, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	role, _ := c.Get("Role")
	if problemSet.GroupId == 0 {
		if role != global.ADMIN && problemSet.UserId != c.GetInt("UserId") {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	} else {
		sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, problemSet.GroupId, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		if role != global.ADMIN && count == 0 {
			c.String(http.StatusForbidden, "没有权限")
			return
		}
	}
	sqlString = `DELETE FROM problem_set WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

// GetWrongCountOfProblemSet godoc
// @Schemes http
// @Description 返回当前用户的该id题库中的错题数量
// @Tags ProblemSet
// @Param id query int true "题库ID"
// @Success 200 {string} string "错题数量"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/statistic/wrong_count [get]
// @Security ApiKeyAuth
func GetWrongCountOfProblemSet(c *gin.Context) {
	sqlString := `SELECT count(*) FROM user_wrong_record WHERE user_id = $1 AND problem_id IN (SELECT problem_id FROM problem_in_problem_set WHERE problem_set_id = $2)`
	var count int
	if err := global.Database.Get(&count, sqlString, c.GetInt("UserId"), c.Query("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, strconv.Itoa(count))
}

// GetFavoriteCountOfProblemSet godoc
// @Schemes http
// @Description 返回当前用户的该id题库中的收藏题目数量
// @Tags ProblemSet
// @Param id query int true "题库ID"
// @Success 200 {string} string "收藏题目数量"
// @Failure default {string} string "服务器错误"
// @Router /problem_set/statistic/fav_count [get]
// @Security ApiKeyAuth
func GetFavoriteCountOfProblemSet(c *gin.Context) {
	sqlString := `SELECT count(*) FROM user_favorite_problem WHERE user_id = $1 AND problem_id IN (SELECT problem_id FROM problem_in_problem_set WHERE problem_set_id = $2)`
	var count int
	if err := global.Database.Get(&count, sqlString, c.GetInt("UserId"), c.Query("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, strconv.Itoa(count))
}
