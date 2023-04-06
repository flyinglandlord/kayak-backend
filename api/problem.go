package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"strconv"
	"time"
)

const (
	ChoiceProblemType = iota
	BlankProblemType
)

func DeleteProblem(c *gin.Context) {
	userId := c.GetInt("UserId")
	problemId := c.Param("id")
	sqlString := `SELECT user_id FROM problem_type WHERE id = $1`
	var problemUserId int
	if err := global.Database.Get(&problemUserId, sqlString, problemId); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if role, _ := c.Get("Role"); userId != problemUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM problem_type WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, problemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

type ProblemFilter struct {
	ID         *int  `json:"id" form:"id"`
	UserId     *int  `json:"user_id" form:"user_id"`
	IsFavorite *bool `json:"is_favorite" form:"is_favorite"`
}
type ChoiceProblemResponse struct {
	ID            int       `json:"id"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserId        int       `json:"user_id"`
	IsPublic      bool      `json:"is_public"`
	IsMultiple    bool      `json:"is_multiple"`
	IsFavorite    bool      `json:"is_favorite"`
	FavoriteCount int       `json:"favorite_count"`
	Choices       []Choice  `json:"choices"`
}
type Choice struct {
	Choice      string `json:"choice"`
	Description string `json:"description"`
}
type AllChoiceProblemResponse struct {
	TotalCount int                     `json:"total_count"`
	Problems   []ChoiceProblemResponse `json:"problems"`
}
type ChoiceProblemCreateRequest struct {
	Description string          `json:"description"`
	IsPublic    bool            `json:"is_public"`
	Choices     []ChoiceRequest `json:"choices"`
}
type ChoiceProblemUpdateRequest struct {
	ID          int             `json:"id"`
	Description *string         `json:"description"`
	IsPublic    *bool           `json:"is_public"`
	Choices     []ChoiceRequest `json:"choices"`
}
type ChoiceRequest struct {
	Choice      string `json:"choice"`
	Description string `json:"description"`
	IsCorrect   bool   `json:"is_correct"`
}

// GetChoiceProblems godoc
// @Schemes http
// @Description 获取符合filter要求的当前用户视角下的所有选择题
// @Param filter query ProblemFilter false "筛选条件"
// @Success 200 {object} []ChoiceProblemResponse "选择题列表"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice/all [get]
// @Security ApiKeyAuth
func GetChoiceProblems(c *gin.Context) {
	sqlString := `SELECT * FROM problem_type` + ` WHERE problem_type_id = ` + strconv.Itoa(ChoiceProblemType)
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString += ` AND is_public = true`
	} else if role == global.USER {
		sqlString += ` AND (is_public = true OR user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
	}
	var filter ProblemFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	if filter.ID != nil {
		sqlString += ` AND id = ` + strconv.Itoa(*filter.ID)
	}
	if filter.UserId != nil {
		sqlString += ` AND user_id = ` + strconv.Itoa(*filter.UserId)
	}
	if filter.IsFavorite != nil {
		if *filter.IsFavorite {
			sqlString += ` AND id IN (SELECT problem_id FROM user_favorite_problem WHERE user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
		} else {
			sqlString += ` AND id NOT IN (SELECT problem_id FROM user_favorite_problem WHERE user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
		}
	}
	var choiceProblems []model.ProblemType
	if err := global.Database.Select(&choiceProblems, sqlString); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var choiceProblemResponses []ChoiceProblemResponse
	for _, problem := range choiceProblems {
		var problemChoices []model.ProblemChoice
		sqlString = `SELECT * FROM problem_choice WHERE id = $1`
		if err := global.Database.Select(&problemChoices, sqlString, problem.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		var choices []Choice
		for _, choice := range problemChoices {
			choices = append(choices, Choice{
				Choice:      choice.Choice,
				Description: choice.Description,
			})
		}
		var CorrectChoiceCount int
		sqlString = `SELECT COUNT(*) FROM problem_choice WHERE id = $1 AND is_correct = true`
		if err := global.Database.Get(&CorrectChoiceCount, sqlString, problem.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		var isFavorite int
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem WHERE user_id = $1 AND problem_id = $2`
		if err := global.Database.Get(&isFavorite, sqlString, c.GetInt("UserId"), problem.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		var favoriteCount int
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem WHERE problem_id = $1`
		if err := global.Database.Get(&favoriteCount, sqlString, problem.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		choiceProblemResponses = append(choiceProblemResponses, ChoiceProblemResponse{
			ID:            problem.ID,
			Description:   problem.Description,
			CreatedAt:     problem.CreatedAt,
			UpdatedAt:     problem.UpdatedAt,
			UserId:        problem.UserId,
			IsPublic:      problem.IsPublic,
			Choices:       choices,
			IsMultiple:    CorrectChoiceCount > 1,
			IsFavorite:    isFavorite > 0,
			FavoriteCount: favoriteCount,
		})
	}
	c.JSON(http.StatusOK, AllChoiceProblemResponse{
		TotalCount: len(choiceProblemResponses),
		Problems:   choiceProblemResponses,
	})
}

// CreateChoiceProblem godoc
// @Schemes http
// @Description 创建选择题
// @Param problem body ChoiceProblemCreateRequest true "选择题信息"
// @Success 200 {object} ChoiceProblemResponse "选择题信息"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice/create [post]
// @Security ApiKeyAuth
func CreateChoiceProblem(c *gin.Context) {
	var request ChoiceProblemCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	tx := global.Database.MustBegin()
	var problemId int
	sqlString := `INSERT INTO problem_type (description, user_id, problem_type_id, is_public, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	if err := global.Database.Get(&problemId, sqlString, request.Description, c.GetInt("UserId"),
		ChoiceProblemType, request.IsPublic, time.Now().Local(), time.Now().Local()); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	for _, choice := range request.Choices {
		sqlString = `INSERT INTO problem_choice (id, choice, description, is_correct) VALUES ($1, $2, $3, $4)`
		if _, err := tx.Exec(sqlString, problemId, choice.Choice, choice.Description, choice.IsCorrect); err != nil {
			_ = tx.Rollback()
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
	}
	if err := tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}

	sqlString = `SELECT * FROM problem_type WHERE id = $1`
	var problem model.ProblemType
	if err := global.Database.Get(&problem, sqlString, problemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var problemChoices []model.ProblemChoice
	sqlString = `SELECT * FROM problem_choice WHERE id = $1`
	if err := global.Database.Select(&problemChoices, sqlString, problemId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var choices []Choice
	for _, choice := range problemChoices {
		choices = append(choices, Choice{
			Choice:      choice.Choice,
			Description: choice.Description,
		})
	}
	var CorrectChoiceCount int
	sqlString = `SELECT COUNT(*) FROM problem_choice WHERE id = $1 AND is_correct = true`
	if err := global.Database.Get(&CorrectChoiceCount, sqlString, problem.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, ChoiceProblemResponse{
		ID:            problem.ID,
		Description:   problem.Description,
		CreatedAt:     problem.CreatedAt,
		UpdatedAt:     problem.UpdatedAt,
		UserId:        problem.UserId,
		IsPublic:      problem.IsPublic,
		IsMultiple:    CorrectChoiceCount > 1,
		IsFavorite:    false,
		FavoriteCount: 0,
		Choices:       choices,
	})
}

// UpdateChoiceProblem godoc
// @Schemes http
// @Description 更新选择题（只需传需要修改的字段,传原值也行）(只有管理员和题目创建者可以更新题目)
// @Param problem body ChoiceProblemUpdateRequest true "选择题信息"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "选择题不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice/update [put]
// @Security ApiKeyAuth
func UpdateChoiceProblem(c *gin.Context) {
	var request ChoiceProblemUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT * FROM problem_type WHERE id = $1`
	var choiceProblem model.ProblemType
	if err := global.Database.Get(&choiceProblem, sqlString, request.ID); err != nil {
		c.String(http.StatusNotFound, "选择题不存在")
		return
	}
	if role, _ := c.Get("Role"); choiceProblem.UserId != c.GetInt("UserId") && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	tx := global.Database.MustBegin()
	if request.Description == nil {
		request.Description = &choiceProblem.Description
	}
	if request.IsPublic == nil {
		request.IsPublic = &choiceProblem.IsPublic
	}
	sqlString = `UPDATE problem_type SET description = $1, is_public = $2, updated_at = $3 WHERE id = $4`
	if _, err := global.Database.Exec(sqlString, request.Description,
		request.IsPublic, time.Now().Local(), request.ID); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	for _, choice := range request.Choices {
		sqlString = `INSERT INTO problem_choice (id, choice, description, is_correct) VALUES ($1, $2, $3, $4) 
			ON CONFLICT (id, choice) DO UPDATE SET description = $3, is_correct = $4`
		if _, err := global.Database.Exec(sqlString, request.ID, choice.Choice,
			choice.Description, choice.IsCorrect); err != nil {
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

type ChoiceProblemAnswerResponse struct {
	Choice      string `json:"choice" db:"choice"`
	Description string `json:"description" db:"description"`
	IsCorrect   bool   `json:"is_correct" db:"is_correct"`
}

// GetChoiceProblemAnswer godoc
// @Schemes http
// @Description 获取选择题答案
// @Param id path int true "选择题ID"
// @Success 200 {object} []ChoiceProblemAnswerResponse "答案信息"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/choice/answer/{id} [get]
// @Security ApiKeyAuth
func GetChoiceProblemAnswer(c *gin.Context) {
	sqlString := `SELECT * FROM problem_type WHERE id = $1 AND problem_type_id = $2`
	var choiceProblem model.ProblemType
	if err := global.Database.Get(&choiceProblem, sqlString, c.Param("id"), ChoiceProblemType); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	if role, _ := c.Get("Role"); choiceProblem.UserId != c.GetInt("UserId") && role != global.ADMIN && !choiceProblem.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT * FROM problem_choice WHERE id = $1`
	var choices []model.ProblemChoice
	if err := global.Database.Select(&choices, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var choiceProblemAnswerResponses []ChoiceProblemAnswerResponse
	for _, choice := range choices {
		choiceProblemAnswerResponses = append(choiceProblemAnswerResponses, ChoiceProblemAnswerResponse{
			Choice:      choice.Choice,
			Description: choice.Description,
			IsCorrect:   choice.IsCorrect,
		})
	}
	c.JSON(http.StatusOK, choiceProblemAnswerResponses)
}

type BlankProblemResponse struct {
	ID            int       `json:"id"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserId        int       `json:"user_id"`
	IsPublic      bool      `json:"is_public"`
	IsFavorite    bool      `json:"is_favorite"`
	FavoriteCount int       `json:"favorite_count"`
}
type AllBlankProblemResponse struct {
	TotalCount int                    `json:"total_count"`
	Problems   []BlankProblemResponse `json:"problems"`
}
type BlankProblemCreateRequest struct {
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	Answer      string `json:"answer"`
}
type BlankProblemUpdateRequest struct {
	ID          int     `json:"id"`
	Description *string `json:"description"`
	IsPublic    *bool   `json:"is_public"`
	Answer      *string `json:"answer"`
}

// GetBlankProblems godoc
// @Schemes http
// @Description 获取当前用户视角下的所有填空题
// @Param filter query ProblemFilter false "筛选条件"
// @Success 200 {object} AllBlankProblemResponse "填空题信息"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /problem/blank/all [get]
// @Security ApiKeyAuth
func GetBlankProblems(c *gin.Context) {
	sqlString := `SELECT * FROM problem_type` + ` WHERE problem_type_id = ` + strconv.Itoa(BlankProblemType)
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString += ` AND is_public = true`
	} else if role == global.USER {
		sqlString += ` AND (is_public = true OR user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
	}
	var filter ProblemFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	if filter.ID != nil {
		sqlString += ` AND id = ` + strconv.Itoa(*filter.ID)
	}
	if filter.UserId != nil {
		sqlString += ` AND user_id = ` + strconv.Itoa(*filter.UserId)
	}
	if filter.IsFavorite != nil {
		if *filter.IsFavorite {
			sqlString += ` AND id IN (SELECT problem_id FROM user_favorite_problem WHERE user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
		} else {
			sqlString += ` AND id NOT IN (SELECT problem_id FROM user_favorite_problem WHERE user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
		}
	}
	var blankProblems []model.ProblemType
	if err := global.Database.Select(&blankProblems, sqlString); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var blankProblemResponses []BlankProblemResponse
	for _, blankProblem := range blankProblems {
		var isFavorite int
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem WHERE user_id = $1 AND problem_id = $2`
		if err := global.Database.Get(&isFavorite, sqlString, c.GetInt("UserId"), blankProblem.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		var favoriteCount int
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem WHERE problem_id = $1`
		if err := global.Database.Get(&favoriteCount, sqlString, blankProblem.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		blankProblemResponses = append(blankProblemResponses, BlankProblemResponse{
			ID:            blankProblem.ID,
			Description:   blankProblem.Description,
			CreatedAt:     blankProblem.CreatedAt,
			UpdatedAt:     blankProblem.UpdatedAt,
			UserId:        blankProblem.UserId,
			IsPublic:      blankProblem.IsPublic,
			IsFavorite:    isFavorite > 0,
			FavoriteCount: favoriteCount,
		})
	}
	c.JSON(http.StatusOK, AllBlankProblemResponse{
		TotalCount: len(blankProblemResponses),
		Problems:   blankProblemResponses,
	})
}

// CreateBlankProblem godoc
// @Schemes http
// @Description 创建填空题
// @Param problem body BlankProblemCreateRequest true "填空题信息"
// @Success 200 {object} BlankProblemResponse "创建成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /problem/blank/create [post]
// @Security ApiKeyAuth
func CreateBlankProblem(c *gin.Context) {
	var request BlankProblemCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	tx := global.Database.MustBegin()
	var problemId int
	sqlString := `INSERT INTO problem_type (problem_type_id, description, is_public, user_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	if err := global.Database.Get(&problemId, sqlString, BlankProblemType, request.Description,
		request.IsPublic, c.GetInt("UserId"), time.Now().Local(), time.Now().Local()); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `INSERT INTO problem_answer (id, answer) VALUES ($1, $2)`
	if _, err := global.Database.Exec(sqlString, problemId, request.Answer); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var problem model.ProblemType
	sqlString = `SELECT * FROM problem_type WHERE id = $1`
	if err := global.Database.Get(&problem, sqlString, problemId); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, BlankProblemResponse{
		ID:            problem.ID,
		Description:   problem.Description,
		CreatedAt:     problem.CreatedAt,
		UpdatedAt:     problem.UpdatedAt,
		UserId:        problem.UserId,
		IsPublic:      problem.IsPublic,
		IsFavorite:    false,
		FavoriteCount: 0,
	})
}

// UpdateBlankProblem godoc
// @Schemes http
// @Description 更新填空题（只有管理员和题目创建者可以更新题目）
// @Param problem body BlankProblemUpdateRequest true "填空题信息"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "填空题不存在"/"答案不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/blank/update [put]
// @Security ApiKeyAuth
func UpdateBlankProblem(c *gin.Context) {
	var request BlankProblemUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT * FROM problem_type WHERE id = $1`
	var blankProblem model.ProblemType
	if err := global.Database.Get(&blankProblem, sqlString, request.ID); err != nil {
		c.String(http.StatusNotFound, "填空题不存在")
		return
	}
	if role, _ := c.Get("Role"); c.GetInt("UserId") != blankProblem.UserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	if request.Description == nil {
		request.Description = &blankProblem.Description
	}
	if request.IsPublic == nil {
		request.IsPublic = &blankProblem.IsPublic
	}
	if request.Answer == nil {
		var answer model.ProblemAnswer
		sqlString = `SELECT answer FROM problem_answer WHERE id = $1`
		if err := global.Database.Get(&answer, sqlString, request.ID); err != nil {
			c.String(http.StatusNotFound, "答案不存在")
			return
		}
		request.Answer = &answer.Answer
	}
	tx := global.Database.MustBegin()
	sqlString = `UPDATE problem_type SET description = $1, is_public = $2, updated_at = $3 WHERE id = $4`
	if _, err := global.Database.Exec(sqlString, request.Description,
		request.IsPublic, time.Now().Local(), request.ID); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `UPDATE problem_answer SET answer = $1 WHERE id = $2`
	if _, err := global.Database.Exec(sqlString, request.Answer, request.ID); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
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

// GetBlankProblemAnswer godoc
// @Schemes http
// @Description 获取填空题答案
// @Param id path int true "填空题ID"
// @Success 200 {string} string "答案"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "填空题不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/blank/answer/{id} [get]
// @Security ApiKeyAuth
func GetBlankProblemAnswer(c *gin.Context) {
	var problem model.ProblemType
	sqlString := `SELECT * FROM problem_type WHERE id = $1`
	if err := global.Database.Get(&problem, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "填空题不存在")
		return
	}
	if role, _ := c.Get("Role"); c.GetInt("UserId") != problem.UserId && role != global.ADMIN && !problem.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT answer FROM problem_answer WHERE id = $1`
	var answer string
	if err := global.Database.Get(&answer, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "填空题不存在")
		return
	}
	c.String(http.StatusOK, answer)
}
