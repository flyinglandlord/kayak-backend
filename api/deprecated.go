package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

// GetUserNotes godoc
// @Schemes http
// @Description 获取当前登录用户的所有笔记
// @Tags Deprecated
// @Success 200 {object} UserNotesResponse "笔记列表"
// @Failure default {string} string "服务器错误"
// @Router /user/note [get]
// @Security ApiKeyAuth
// @Deprecated
func GetUserNotes(c *gin.Context) {
	var notes []model.Note
	sqlString := `SELECT * FROM note WHERE user_id = $1`
	if err := global.Database.Select(&notes, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var noteResponses []NoteResponse
	for _, note := range notes {
		var likeCount, favoriteCount int
		sqlString = `SELECT COUNT(*) FROM user_like_note WHERE note_id = $1 AND user_id = $2`
		if err := global.Database.Get(&likeCount, sqlString, note.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		sqlString = `SELECT COUNT(*) FROM user_favorite_note WHERE note_id = $1 AND user_id = $2`
		if err := global.Database.Get(&favoriteCount, sqlString, note.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		noteResponses = append(noteResponses, NoteResponse{
			ID:            note.ID,
			Title:         note.Title,
			Content:       note.Content,
			CreatedAt:     note.CreatedAt,
			IsLiked:       likeCount > 0,
			IsFavorite:    favoriteCount > 0,
			LikeCount:     likeCount,
			FavoriteCount: favoriteCount,
		})
	}
	c.JSON(http.StatusOK, UserNotesResponse{
		TotalCount: len(notes),
		Notes:      noteResponses,
	})
}

type FavoriteProblemResponse struct {
	TotalCount int                 `json:"total_count"`
	Problems   []model.ProblemType `json:"problems"`
}

// GetUserFavoriteProblems godoc
// @Schemes http
// @Description 获取当前登录用户收藏的题目
// @Tags Deprecated
// @Success 200 {object} FavoriteProblemResponse "收藏的题目列表"
// @Failure default {string} string "服务器错误"
// @Router /user/favorite/problem [get]
// @Security ApiKeyAuth
// @Deprecated
func GetUserFavoriteProblems(c *gin.Context) {
	var problems []model.ProblemType
	userId := c.GetInt("UserId")
	sqlString :=
		`SELECT p.id AS id, p.description as description, p.created_at AS created_at, 
    		p.updated_at AS updated_at, p.user_id AS user_id, p.problem_type_id AS problem_type_id, p.is_public AS is_public
	     FROM user_favorite_problem ufp JOIN problem_type p ON ufp.problem_id = p.id 
	     WHERE ufp.user_id = $1`
	if err := global.Database.Select(&problems, sqlString, userId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, FavoriteProblemResponse{
		TotalCount: len(problems),
		Problems:   problems,
	})
}

type FavoriteProblemSetResponse struct {
	TotalCount  int                  `json:"total_count"`
	ProblemSets []ProblemSetResponse `json:"problem_set"`
}

// GetUserFavoriteProblemSets godoc
// @Schemes http
// @Description 获取当前登录用户收藏的题集
// @Tags Deprecated
// @Success 200 {object} FavoriteProblemSetResponse "收藏的题集列表"
// @Failure default {string} string "服务器错误"
// @Router /user/favorite/problem_set [get]
// @Security ApiKeyAuth
// @Deprecated
func GetUserFavoriteProblemSets(c *gin.Context) {
	var problemsets []model.ProblemSet
	userId := c.GetInt("UserId")
	sqlString :=
		`SELECT ps.id AS id, ps.name AS name, ps.description AS description, ps.created_at AS created_at, 
    		ps.updated_at AS updated_at, ps.user_id AS user_id, ps.is_public AS is_public
	 	 FROM user_favorite_problem_set ufps RIGHT JOIN problem_set ps ON ufps."problem_set_id" = ps.id
	 	 WHERE ufps.user_id = $1`
	if err := global.Database.Select(&problemsets, sqlString, userId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var problemSetResponses []ProblemSetResponse
	for _, problemset := range problemsets {
		// 获取收藏数
		var favoriteCount int
		sqlString := `SELECT count(*) FROM user_favorite_problem_set WHERE "problem_set_id" = $1`
		if err := global.Database.Get(&favoriteCount, sqlString, problemset.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		// 获取题目数
		var problemCount int
		sqlString = `SELECT count(*) FROM problem_in_problem_set WHERE "problem_set_id" = $1`
		if err := global.Database.Get(&problemCount, sqlString, problemset.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		// 获取自己是否点过赞
		var isFavorite bool
		sqlString = `SELECT count(*) FROM user_favorite_problem_set WHERE "problem_set_id" = $1 AND user_id = $2`
		if err := global.Database.Get(&isFavorite, sqlString, problemset.ID, userId); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		problemSetResponses = append(problemSetResponses, ProblemSetResponse{
			ID:            problemset.ID,
			Name:          problemset.Name,
			Description:   problemset.Description,
			CreatedAt:     problemset.CreatedAt,
			ProblemCount:  problemCount,
			FavoriteCount: favoriteCount,
			IsFavorite:    isFavorite,
			UserId:        problemset.UserId,
			IsPublic:      problemset.IsPublic,
		})
	}
	c.JSON(http.StatusOK, FavoriteProblemSetResponse{
		TotalCount:  len(problemSetResponses),
		ProblemSets: problemSetResponses,
	})
}

type FavoriteNoteResponse struct {
	TotalCount int            `json:"total_count"`
	Notes      []NoteResponse `json:"notes"`
}

// GetUserFavoriteNotes godoc
// @Schemes http
// @Description 获取当前登录用户收藏的笔记
// @Tags Deprecated
// @Success 200 {object} FavoriteNoteResponse "收藏的笔记列表"
// @Failure default {string} string "服务器错误"
// @Router /user/favorite/note [get]
// @Security ApiKeyAuth
// @Deprecated
func GetUserFavoriteNotes(c *gin.Context) {
	var notes []model.Note
	userId := c.GetInt("UserId")
	sqlString :=
		`SELECT n.id AS id, n.title AS title, n.content AS content, 
       		n.created_at AS created_at, n.user_id AS user_id, n.updated_at AS updated_at, n.is_public AS is_public
		 FROM user_favorite_note ufn JOIN note n ON ufn.note_id = n.id 
		 WHERE ufn.user_id = $1`
	if err := global.Database.Select(&notes, sqlString, userId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var noteResponses []NoteResponse
	for _, note := range notes {
		var likeCount, favoriteCount int
		sqlString = `SELECT COUNT(*) FROM user_like_note WHERE note_id = $1`
		if err := global.Database.Get(&likeCount, sqlString, note.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		sqlString = `SELECT COUNT(*) FROM user_favorite_note WHERE note_id = $1`
		if err := global.Database.Get(&favoriteCount, sqlString, note.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		noteResponses = append(noteResponses, NoteResponse{
			ID:            note.ID,
			Title:         note.Title,
			Content:       note.Content,
			CreatedAt:     note.CreatedAt,
			UserId:        note.UserId,
			IsPublic:      note.IsPublic,
			LikeCount:     likeCount,
			FavoriteCount: favoriteCount,
			IsFavorite:    favoriteCount > 0,
			IsLiked:       likeCount > 0,
		})
	}
	c.JSON(http.StatusOK, FavoriteNoteResponse{
		TotalCount: len(noteResponses),
		Notes:      noteResponses,
	})
}

// GetUserProblemSets godoc
// @Schemes http
// @Description 获取当前登录用户的所有题集
// @Tags Deprecated
// @Success 200 {object} []ProblemSetResponse "题集列表"
// @Failure default {string} string "服务器错误"
// @Router /user/problem_set [get]
// @Security ApiKeyAuth
// @Deprecated
func GetUserProblemSets(c *gin.Context) {
	var problemsets []model.ProblemSet
	userId := c.GetInt("UserId")
	sqlString :=
		`SELECT ps.id AS id, ps.name AS name, ps.description AS description, ps.created_at AS created_at, 
    		ps.updated_at AS updated_at, ps.user_id AS user_id, ps.is_public AS is_public
	 	 FROM problem_set ps
	 	 WHERE ps.user_id = $1`
	if err := global.Database.Select(&problemsets, sqlString, userId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var problemsetResponses []ProblemSetResponse
	for _, problemset := range problemsets {
		// 获取收藏数
		var favoriteCount int
		sqlString := `SELECT count(*) FROM user_favorite_problem_set WHERE "problem_set_id" = $1`
		if err := global.Database.Get(&favoriteCount, sqlString, problemset.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		// 获取题目数
		var problemCount int
		sqlString = `SELECT count(*) FROM problem_in_problem_set WHERE "problem_set_id" = $1`
		if err := global.Database.Get(&problemCount, sqlString, problemset.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		// 获取自己是否点过赞
		var isFavorite bool
		sqlString = `SELECT count(*) FROM user_favorite_problem_set WHERE "problem_set_id" = $1 AND user_id = $2`
		if err := global.Database.Get(&isFavorite, sqlString, problemset.ID, userId); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		problemsetResponses = append(problemsetResponses, ProblemSetResponse{
			ID:            problemset.ID,
			Name:          problemset.Name,
			Description:   problemset.Description,
			CreatedAt:     problemset.CreatedAt,
			UpdatedAt:     problemset.UpdatedAt,
			UserId:        problemset.UserId,
			ProblemCount:  problemCount,
			FavoriteCount: favoriteCount,
			IsFavorite:    isFavorite,
			IsPublic:      problemset.IsPublic,
		})
	}
	c.JSON(http.StatusOK, problemsetResponses)
}

type ChoiceProblemItem struct {
	Id            int    `json:"id"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	UserId        int    `json:"user_id"`
	IsPublic      bool   `json:"is_public"`
	IsMultiple    bool   `json:"is_multiple"`
	FavoriteCount int    `json:"favorite_count"`
}

// GetUserChoiceProblems godoc
// @Schemes http
// @Description 获取当前登录用户的所有选择题
// @Tags Deprecated
// @Success 200 {object} []ChoiceProblemResponse "选择题列表"
// @Failure default {string} string "服务器错误"
// @Router /user/problem/choice [get]
// @Security ApiKeyAuth
// @Deprecated
func GetUserChoiceProblems(c *gin.Context) {
	var problemType []model.ProblemType
	var choiceProblems []ChoiceProblemItem
	userId := c.GetInt("UserId")
	sqlString := `SELECT id, description, created_at, updated_at, user_id, is_public FROM problem_type WHERE problem_type_id = $1 AND user_id = $2`
	if err := global.Database.Select(&problemType, sqlString, ChoiceProblemType, userId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	for _, item := range problemType {
		choiceProblem := ChoiceProblemItem{
			Id:            item.ID,
			Description:   item.Description,
			CreatedAt:     item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     item.UpdatedAt.Format(time.RFC3339),
			UserId:        item.UserId,
			IsPublic:      item.IsPublic,
			FavoriteCount: 0,
		}
		var CorrectAnswerCount int
		sqlString := `SELECT COUNT(*) FROM problem_choice WHERE id = $1 AND is_correct = true`
		if err := global.Database.Get(&CorrectAnswerCount, sqlString, item.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		choiceProblem.IsMultiple = CorrectAnswerCount > 1
		// 获取收藏数
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem WHERE problem_id = $1`
		if err := global.Database.Get(&choiceProblem.FavoriteCount, sqlString, item.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		choiceProblems = append(choiceProblems, choiceProblem)
	}

	c.JSON(http.StatusOK, choiceProblems)
}

// GetUserBlankProblems godoc
// @Schemes http
// @Description 获取当前登录用户的所有填空题
// @Tags Deprecated
// @Success 200 {object} BlankProblemResponse "填空题信息"
// @Failure default {string} string "服务器错误"
// @Router /user/problem/blank [get]
// @Security ApiKeyAuth
// @Deprecated
func GetUserBlankProblems(c *gin.Context) {
	userId := c.GetInt("UserId")
	sqlString := `SELECT id, description FROM problem_type WHERE problem_type_id = $1 AND user_id = $2`
	var blankProblems []BlankProblemResponse
	if err := global.Database.Select(&blankProblems, sqlString, BlankProblemType, userId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	for _, blankProblem := range blankProblems {
		sqlString := `SELECT count(*) FROM user_favorite_problem WHERE problem_id = $1`
		if err := global.Database.Get(&blankProblem.FavoriteCount, sqlString, blankProblem.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
	}
	c.JSON(http.StatusOK, blankProblems)
}

// GetChoiceProblem godoc
// @Schemes http
// @Description 获取单个选择题信息（只有管理员和题目创建者可以获取私有题目）
// @Tags Deprecated
// @Param id path int true "选择题ID"
// @Success 200 {object} ChoiceProblemResponse "选择题信息"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "选择题不存在"
// @Failure default {string} string "获取选项失败"
// @Router /problem/choice/{id} [get]
// @Security ApiKeyAuth
// @Deprecated
func GetChoiceProblem(c *gin.Context) {
	var choiceProblem model.ProblemType
	sqlString := `SELECT * FROM problem_type WHERE problem_type_id = $1 AND id = $2`
	if err := global.Database.Get(&choiceProblem, sqlString, ChoiceProblemType, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "选择题不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && choiceProblem.UserId != c.GetInt("UserId") && !choiceProblem.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	var problemChoices []model.ProblemChoice
	sqlString = `SELECT * FROM problem_choice WHERE id = $1`
	if err := global.Database.Select(&problemChoices, sqlString, choiceProblem.ID); err != nil {
		c.String(http.StatusInternalServerError, "获取选项失败")
		return
	}
	var choices []Choice
	for _, choice := range problemChoices {
		choices = append(choices, Choice{
			Choice:      choice.Choice,
			Description: choice.Description,
		})
	}

	// 获取收藏数
	var favoriteCount int
	sqlString = `SELECT COUNT(*) FROM user_favorite_problem WHERE problem_id = $1`
	if err := global.Database.Get(&favoriteCount, sqlString, choiceProblem.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}

	choiceProblemResponse := ChoiceProblemResponse{
		ID:            choiceProblem.ID,
		Description:   choiceProblem.Description,
		CreatedAt:     choiceProblem.CreatedAt,
		UpdatedAt:     choiceProblem.UpdatedAt,
		UserId:        choiceProblem.UserId,
		IsPublic:      choiceProblem.IsPublic,
		Choices:       choices,
		FavoriteCount: favoriteCount,
	}
	c.JSON(http.StatusOK, choiceProblemResponse)
}

// GetProblemSetContainsProblem godoc
// @Schemes http
// @Description 获取包含某题目的题集
// @Tags Deprecated
// @Param id path int true "题目ID"
// @Success 200 {array} ProblemSetResponse "题集列表"
// @Failure 404 {string} string "题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/{id}/problem_set [get]
// @Security ApiKeyAuth
// @Deprecated
func GetProblemSetContainsProblem(c *gin.Context) {
	var problemSetList []model.ProblemSet
	sqlString := `SELECT * FROM problem_set WHERE id IN (SELECT problem_set_id FROM problem_in_problem_set WHERE problem_id = $1)`
	if err := global.Database.Select(&problemSetList, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	var problemSetResponseList []ProblemSetResponse
	for _, problemSet := range problemSetList {
		var ProblemCount int
		sqlString = `SELECT COUNT(*) FROM problem_in_problem_set WHERE problem_set_id = $1`
		if err := global.Database.Get(&ProblemCount, sqlString, problemSet.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		var isFavorite bool
		sqlString = `SELECT * FROM user_favorite_problem_set WHERE user_id = $1 AND problem_set_id = $2`
		if err := global.Database.Get(&isFavorite, sqlString, c.GetInt("UserId"), problemSet.ID); err != nil {
			isFavorite = false
		}
		problemSetResponseList = append(problemSetResponseList, ProblemSetResponse{
			ID:           problemSet.ID,
			Name:         problemSet.Name,
			Description:  problemSet.Description,
			UserId:       problemSet.UserId,
			CreatedAt:    problemSet.CreatedAt,
			UpdatedAt:    problemSet.UpdatedAt,
			ProblemCount: ProblemCount,
			IsFavorite:   isFavorite,
		})
	}
	c.JSON(http.StatusOK, problemSetResponseList)
}

// GetBlankProblem godoc
// @Schemes http
// @Description 获取单个填空题信息（只有管理员和题目创建者可以查看私有题目）
// @Tags Deprecated
// @Param id path int true "填空题ID"
// @Success 200 {object} BlankProblemResponse "填空题信息"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "填空题不存在"
// @Failure default {string} string "服务器错误"
// @Router /problem/blank/{id} [get]
// @Deprecated
func GetBlankProblem(c *gin.Context) {
	var blankProblem model.ProblemType
	sqlString := `SELECT * FROM problem_type WHERE problem_type_id = $1 AND id = $2`
	if err := global.Database.Get(&blankProblem, sqlString, BlankProblemType, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "填空题不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && blankProblem.UserId != c.GetInt("UserId") && !blankProblem.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	// 获取收藏数
	var favoriteCount int
	sqlString = `SELECT COUNT(*) FROM user_favorite_problem WHERE problem_id = $1`
	if err := global.Database.Get(&favoriteCount, sqlString, blankProblem.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	blankProblemResponse := BlankProblemResponse{
		ID:            blankProblem.ID,
		Description:   blankProblem.Description,
		CreatedAt:     blankProblem.CreatedAt,
		UpdatedAt:     blankProblem.UpdatedAt,
		UserId:        blankProblem.UserId,
		IsPublic:      blankProblem.IsPublic,
		FavoriteCount: favoriteCount,
	}
	c.JSON(http.StatusOK, blankProblemResponse)
}
