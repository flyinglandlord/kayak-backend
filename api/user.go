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
	c.JSON(200, userInfo)
}

type UserInfo struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

// UpdateUserInfo godoc
// @Schemes http
// @Description 更新用户信息
// @Param info body UserInfo true "用户信息"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /user/update [put]
// @Security ApiKeyAuth
func UpdateUserInfo(c *gin.Context) {
	user := UserInfo{}
	sqlString := `UPDATE "user" SET name = $1, email = $2, phone = $3 WHERE id = $4`
	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, "请求格式错误")
		return
	}
	if _, err := global.Database.Exec(sqlString, user.Name, user.Email, user.Phone, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(200, "更新成功")
}

// GetUserNotes godoc
// @Schemes http
// @Description 获取当前登录用户的所有笔记
// @Success 200 {object} []NoteResponse "笔记列表"
// @Failure default {string} string "服务器错误"
// @Router /user/note [get]
// @Security ApiKeyAuth
func GetUserNotes(c *gin.Context) {
	var notes []NoteResponse
	sqlString := `SELECT id, title, content FROM note WHERE user_id = $1`
	if err := global.Database.Select(&notes, sqlString, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, notes)
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

// GetUserFavoriteProblems godoc
// @Schemes http
// @Description 获取当前登录用户收藏的题目
// @Success 200 {object} []FavoriteProblemResponse "收藏的题目列表"
// @Failure default {string} string "服务器错误"
// @Router /user/favorite/problem [get]
func GetUserFavoriteProblems(c *gin.Context) {
	var problems []FavoriteProblemResponse
	userId := c.GetInt("UserId")
	sqlString := "SELECT problem_id FROM user_favorite_problem WHERE user_id = $1"
	if err := global.Database.Get(&problems, sqlString, userId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, problems)
}

// GetUserFavoriteProblemsets godoc
// @Schemes http
// @Description 获取当前登录用户收藏的题集
// @Success 200 {object} []FavoriteProblemsetResponse "收藏的题集列表"
// @Failure default {string} string "服务器错误"
// @Router /user/favorite/problemset [get]
func GetUserFavoriteProblemsets(c *gin.Context) {
	var problemsets []FavoriteProblemsetResponse
	userId := c.GetInt("UserId")
	sqlString := "SELECT problemset_id FROM user_favorite_problemset WHERE user_id = $1"
	if err := global.Database.Get(&problemsets, sqlString, userId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, problemsets)
}

// GetUserProblemsets godoc
// @Schemes http
// @Description 获取当前登录用户的所有题集
// @Success 200 {object} []ProblemsetResponse "题集列表"
// @Failure default {string} string "服务器错误"
// @Router /user/problemset [get]
// @Security ApiKeyAuth
func GetUserProblemsets(c *gin.Context) {
	var problemsets []ProblemsetResponse
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
