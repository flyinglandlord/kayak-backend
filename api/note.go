package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type NoteResponse struct {
	ID         int       `json:"id" db:"id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	IsLiked    bool      `json:"is_liked"`
	IsFavorite bool      `json:"is_favorite"`
}
type NoteCreateRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	IsPublic bool   `json:"is_public"`
}
type NoteUpdateRequest struct {
	ID       int     `json:"id"`
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	IsPublic *bool   `json:"is_public"`
}

// GetNotes godoc
// @Schemes http
// @Description 获取当前用户视角下的所有笔记
// @Param id query int false "笔记ID"
// @Success 200 {object} []NoteResponse "笔记列表"
// @Failure default {string} string "服务器错误"
// @Router /note/all [get]
// @Security ApiKeyAuth
func GetNotes(c *gin.Context) {
	var notes []model.Note
	var sqlString string
	var err error
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString = `SELECT * FROM note WHERE is_public = true`
		if c.Query("id") != "" {
			sqlString += ` AND id = ` + c.Query("id")
		}
		err = global.Database.Select(&notes, sqlString)
	} else if role == global.USER {
		sqlString = `SELECT * FROM note WHERE (is_public = true OR user_id = $1)`
		if c.Query("id") != "" {
			sqlString += ` AND id = ` + c.Query("id")
		}
		err = global.Database.Select(&notes, sqlString, c.GetInt("UserId"))
	} else {
		sqlString = `SELECT * FROM note`
		if c.Query("id") != "" {
			sqlString += ` WHERE id = ` + c.Query("id")
		}
		err = global.Database.Select(&notes, sqlString)
	}
	if err != nil {
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
	c.JSON(http.StatusOK, noteResponses)
}

// CreateNote godoc
// @Schemes http
// @Description 创建笔记
// @Param note body NoteCreateRequest true "笔记信息"
// @Success 200 {string} string "创建成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /note/create [post]
// @Security ApiKeyAuth
func CreateNote(c *gin.Context) {
	var request NoteCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `INSERT INTO note (title, content, created_at, updated_at, user_id, is_public) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := global.Database.Exec(sqlString, request.Title, request.Content, time.Now().Local(),
		time.Now().Local(), c.GetInt("UserId"), request.IsPublic); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

// UpdateNote godoc
// @Schemes http
// @Description 更新笔记（只有管理员和笔记作者可以更新）(可以只传需要更新的字段)
// @Param note body NoteUpdateRequest true "笔记信息"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/update [put]
// @Security ApiKeyAuth
func UpdateNote(c *gin.Context) {
	var note model.Note
	var request NoteUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT * FROM note WHERE id = $1`
	if err := global.Database.Get(&note, sqlString, request.ID); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); note.UserId != c.GetInt("UserId") && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	if request.Title == nil {
		request.Title = &note.Title
	}
	if request.Content == nil {
		request.Content = &note.Content
	}
	if request.IsPublic == nil {
		request.IsPublic = &note.IsPublic
	}
	sqlString = `UPDATE note SET title = $1, content = $2, updated_at = $3, is_public = $4 WHERE id = $5`
	if _, err := global.Database.Exec(sqlString, request.Title, request.Content,
		time.Now().Local(), request.IsPublic, request.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "更新成功")
}

// DeleteNote godoc
// @Schemes http
// @Description 删除笔记（只有管理员和笔记作者可以删除）
// @Param id path int true "笔记ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteNote(c *gin.Context) {
	userId := c.GetInt("UserId")
	noteId := c.Param("id")
	sqlString := `SELECT user_id FROM note WHERE id = $1`
	var noteUserId int
	if err := global.Database.Get(&noteUserId, sqlString, noteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); userId != noteUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM note WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, noteId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

// LikeNote godoc
// @Schemes http
// @Description 点赞笔记
// @Param id path int true "笔记ID"
// @Success 200 {string} string "点赞成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/like/{id} [post]
// @Security ApiKeyAuth
func LikeNote(c *gin.Context) {
	var note model.Note
	userId := c.GetInt("UserId")
	noteId := c.Param("id")
	sqlString := `SELECT * FROM note WHERE id = $1`
	if err := global.Database.Get(&note, sqlString, noteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && note.UserId != userId && !note.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO user_like_note (user_id, note_id, created_at) VALUES ($1, $2, $3)`
	if _, err := global.Database.Exec(sqlString, userId, noteId, time.Now().Local()); err != nil {
		return
	}
	c.String(http.StatusOK, "点赞成功")
}

// UnlikeNote godoc
// @Schemes http
// @Description 取消点赞笔记
// @Param id path int true "笔记ID"
// @Success 200 {string} string "取消点赞成功"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/unlike/{id} [post]
// @Security ApiKeyAuth
func UnlikeNote(c *gin.Context) {
	var note model.Note
	userId := c.GetInt("UserId")
	noteId := c.Param("id")
	sqlString := `SELECT * FROM note WHERE id = $1`
	if err := global.Database.Get(&note, sqlString, noteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	sqlString = `DELETE FROM user_like_note WHERE user_id = $1 AND note_id = $2`
	if _, err := global.Database.Exec(sqlString, userId, noteId); err != nil {
		return
	}
	c.String(http.StatusOK, "取消点赞成功")
}

// FavoriteNote godoc
// @Schemes http
// @Description 收藏笔记
// @Param id path int true "笔记ID"
// @Success 200 {string} string "收藏成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/favorite/{id} [post]
// @Security ApiKeyAuth
func FavoriteNote(c *gin.Context) {
	var note model.Note
	userId := c.GetInt("UserId")
	noteId := c.Param("id")
	sqlString := `SELECT * FROM note WHERE id = $1`
	if err := global.Database.Get(&note, sqlString, noteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && note.UserId != userId && !note.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO user_favorite_note (user_id, note_id, created_at) VALUES ($1, $2, $3)`
	if _, err := global.Database.Exec(sqlString, userId, noteId, time.Now().Local()); err != nil {
		return
	}
	c.String(http.StatusOK, "收藏成功")
}

// UnfavoriteNote godoc
// @Schemes http
// @Description 取消收藏笔记
// @Param id path int true "笔记ID"
// @Success 200 {string} string "取消收藏成功"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/unfavorite/{id} [post]
// @Security ApiKeyAuth
func UnfavoriteNote(c *gin.Context) {
	var note model.Note
	userId := c.GetInt("UserId")
	noteId := c.Param("id")
	sqlString := `SELECT * FROM note WHERE id = $1`
	if err := global.Database.Get(&note, sqlString, noteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	sqlString = `DELETE FROM user_favorite_note WHERE user_id = $1 AND note_id = $2`
	if _, err := global.Database.Exec(sqlString, userId, noteId); err != nil {
		return
	}
	c.String(http.StatusOK, "取消收藏成功")
}
