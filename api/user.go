package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type UserInfoResponse struct {
	UserId     int       `json:"user_id"`
	UserName   string    `json:"user_name"`
	Email      *string   `json:"email"`
	Phone      *string   `json:"phone"`
	AvatarPath string    `json:"avatar_path"`
	CreateAt   time.Time `json:"create_at"`
}

// GetUserInfo godoc
// @Schemes http
// @Description 获取用户信息
// @Success 200 {object} UserInfoResponse "用户信息"
// @Failure default {string} string "服务器错误"
// @Router /user/info [get]
// @Security ApiKeyAuth
func GetUserInfo(c *gin.Context) {
	user := model.User{}
	sqlString := `SELECT name, email, phone, avatar_url, created_at FROM "user" WHERE id = $1`
	if err := global.Database.Get(&user, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	userInfo := UserInfoResponse{
		UserId:     c.GetInt("UserId"),
		UserName:   user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		AvatarPath: user.AvatarURL,
		CreateAt:   user.CreatedAt,
	}
	c.JSON(http.StatusOK, userInfo)
}

type UserInfoRequest struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

// UpdateUserInfo godoc
// @Schemes http
// @Description 更新用户信息
// @Param info body UserInfoRequest true "用户信息"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /user/update [put]
// @Security ApiKeyAuth
func UpdateUserInfo(c *gin.Context) {
	user := UserInfoRequest{}
	sqlString := `UPDATE "user" SET name = $1, email = $2, phone = $3 WHERE id = $4`
	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, "请求格式错误")
		return
	}
	if _, err := global.Database.Exec(sqlString, user.Name, user.Email, user.Phone, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "更新成功")
}

type UserNotesResponse struct {
	TotalCount int            `json:"total_count"`
	Notes      []NoteResponse `json:"notes"`
}

// GetUserNotes godoc
// @Schemes http
// @Description 获取当前登录用户的所有笔记
// @Success 200 {object} UserNotesResponse "笔记列表"
// @Failure default {string} string "服务器错误"
// @Router /user/note [get]
// @Security ApiKeyAuth
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
			ID:         note.ID,
			Title:      note.Title,
			Content:    note.Content,
			CreatedAt:  note.CreatedAt,
			IsLiked:    likeCount > 0,
			IsFavorite: favoriteCount > 0,
		})
	}
	c.JSON(http.StatusOK, UserNotesResponse{
		TotalCount: len(notes),
		Notes:      noteResponses,
	})
}

// GetUserWrongRecords godoc
// @Schemes http
// @Description 获取当前登录用户的所有错题记录
// @Success 200 {object} []WrongRecordResponse "错题记录列表"
// @Failure default {string} string "服务器错误"
// @Router /user/wrong_record [get]
// @Security ApiKeyAuth
func GetUserWrongRecords(c *gin.Context) {
	var wrongRecords []WrongRecordResponse
	sqlString := `SELECT problem_id, count, created_at, updated_at FROM user_wrong_record WHERE user_id = $1`
	if err := global.Database.Select(&wrongRecords, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, wrongRecords)
}

type FavoriteProblemResponse struct {
	TotalCount int                 `json:"total_count"`
	Problems   []model.ProblemType `json:"problems"`
}

// GetUserFavoriteProblems godoc
// @Schemes http
// @Description 获取当前登录用户收藏的题目
// @Success 200 {object} FavoriteProblemResponse "收藏的题目列表"
// @Failure default {string} string "服务器错误"
// @Router /user/favorite/problem [get]
// @Security ApiKeyAuth
func GetUserFavoriteProblems(c *gin.Context) {
	var problems []model.ProblemType
	userId := c.GetInt("UserId")
	sqlString :=
		`SELECT p.id AS id, p.description as description, p.created_at AS created_at, 
    		p.updated_at AS updated_at, p.user_id AS user_id, p.problem_type_id AS problem_type_id, p.is_public AS is_public
	     FROM user_favorite_problem ufp JOIN problem_type p ON ufp.problem_id = p.id 
	     WHERE ufp.user_id = $1 GROUP BY ufp.problem_id`
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

// GetUserFavoriteProblemsets godoc
// @Schemes http
// @Description 获取当前登录用户收藏的题集
// @Success 200 {object} FavoriteProblemSetResponse "收藏的题集列表"
// @Failure default {string} string "服务器错误"
// @Router /user/favorite/problemset [get]
// @Security ApiKeyAuth
func GetUserFavoriteProblemSets(c *gin.Context) {
	var problemsets []ProblemSetResponse
	userId := c.GetInt("UserId")
	sqlString :=
		`SELECT ps.id AS id, ps.name AS name, ps.description AS description, ps.created_at AS created_at, 
    		ps.updated_at AS updated_at, count(*) AS problem_count
	 	 FROM user_favorite_problemset ufps JOIN problemset ps ON ufps."problemSet_id" = ps.id JOIN problem_in_problemset pip on ps.id = pip."problemSet_id"
	 	 WHERE ufps.user_id = $1 GROUP BY ps.id`
	if err := global.Database.Select(&problemsets, sqlString, userId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, FavoriteProblemSetResponse{
		TotalCount:  len(problemsets),
		ProblemSets: problemsets,
	})
}

type FavoriteNoteResponse struct {
	TotalCount int          `json:"total_count"`
	Notes      []model.Note `json:"notes"`
}

// GetUserFavoriteNotes godoc
// @Schemes http
// @Description 获取当前登录用户收藏的笔记
// @Success 200 {object} FavoriteNoteResponse "收藏的笔记列表"
// @Failure default {string} string "服务器错误"
// @Router /user/favorite/note [get]
// @Security ApiKeyAuth
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
	c.JSON(http.StatusOK, FavoriteNoteResponse{
		TotalCount: len(notes),
		Notes:      notes,
	})
}

// GetUserProblemsets godoc
// @Schemes http
// @Description 获取当前登录用户的所有题集
// @Success 200 {object} []ProblemSetResponse "题集列表"
// @Failure default {string} string "服务器错误"
// @Router /user/problemset [get]
// @Security ApiKeyAuth
func GetUserProblemSets(c *gin.Context) {
	var problemsets []ProblemSetResponse
	sqlString := `SELECT id, name, description FROM problemset WHERE user_id = $1`
	if err := global.Database.Select(&problemsets, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, problemsets)
}

// GetUserChoiceProblems godoc
// @Schemes http
// @Description 获取当前登录用户的所有选择题
// @Success 200 {object} []ChoiceProblemResponse "选择题列表"
// @Failure default {string} string "服务器错误"
// @Router /user/problem/choice [get]
// @Security ApiKeyAuth
func GetUserChoiceProblems(c *gin.Context) {
	var choiceProblems []ChoiceProblemResponse
	sqlString := `SELECT id, description FROM problem_type WHERE problem_type_id = $1 AND user_id = $2`
	if err := global.Database.Select(&choiceProblems, sqlString, ChoiceProblemType, c.GetInt("UserId")); err != nil {
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

// GetUserBlankProblems godoc
// @Schemes http
// @Description 获取当前登录用户的所有填空题
// @Success 200 {object} BlankProblemResponse "填空题信息"
// @Failure default {string} string "服务器错误"
// @Router /user/problem/blank [get]
// @Security ApiKeyAuth
func GetUserBlankProblems(c *gin.Context) {
	userId := c.GetInt("UserId")
	sqlString := `SELECT id, description FROM problem_type WHERE problem_type_id = $1 AND user_id = $2`
	var blankProblems []BlankProblemResponse
	if err := global.Database.Select(&blankProblems, sqlString, BlankProblemType, userId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, blankProblems)
}
