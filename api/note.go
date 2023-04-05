package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type NoteItem struct {
	ID            int    `json:"id" db:"id"`
	Title         string `json:"title" db:"title"`
	Content       string `json:"content" db:"content"`
	CreatedAt     string `json:"created_at" db:"created_at"`
	IsLiked       bool   `json:"is_liked"`
	IsFavorite    bool   `json:"is_favorite"`
	LikeCount     int    `json:"like_count"`
	FavoriteCount int    `json:"favorite_count"`
}

type NoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteResponse struct {
	TotalCount int        `json:"total_count"`
	Notes      []NoteItem `json:"notes"`
}

// GetNotes godoc
// @Schemes http
// @Description 获取当前用户视角下的所有笔记
// @Success 200 {object} []NoteItem "笔记列表"
// @Failure default {string} string "服务器错误"
// @Router /note/all [get]
// @Security ApiKeyAuth
func GetNotes(c *gin.Context) {
	var notes []model.Note
	var sqlString string
	var err error
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString = `SELECT id, title, content, created_at FROM note WHERE is_public = true`
		err = global.Database.Select(&notes, sqlString)
	} else if role == global.USER {
		sqlString = `SELECT id, title, content, created_at FROM note WHERE (is_public = true OR user_id = $1)`
		err = global.Database.Select(&notes, sqlString, c.GetInt("UserId"))
	} else {
		sqlString = `SELECT id, title, content, created_at FROM note`
		err = global.Database.Select(&notes, sqlString)
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var noteResponses []NoteItem
	for _, note := range notes {
		var noteResponse NoteItem
		noteResponse.ID = note.ID
		noteResponse.Title = note.Title
		noteResponse.Content = note.Content
		noteResponse.CreatedAt = note.CreatedAt.Format(time.RFC3339)

		// 查询是否点赞
		sqlString = `SELECT COUNT(*) FROM user_like_note WHERE note_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, note.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		noteResponse.LikeCount = count
		if count > 0 {
			noteResponse.IsLiked = true
		}

		// 查询是否收藏
		sqlString = `SELECT COUNT(*) FROM user_favorite_note WHERE note_id = $1 AND user_id = $2`
		if err := global.Database.Get(&count, sqlString, note.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		noteResponse.FavoriteCount = count
		if count > 0 {
			noteResponse.IsFavorite = true
		}
		noteResponses = append(noteResponses, noteResponse)
	}
	c.JSON(http.StatusOK, NoteResponse{
		TotalCount: len(noteResponses),
		Notes:      noteResponses,
	})
}

// CreateNote godoc
// @Schemes http
// @Description 创建笔记
// @Param note body NoteRequest true "笔记信息"
// @Param is_public query bool true "是否公开"
// @Success 200 {string} string "创建成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /note/create [post]
// @Security ApiKeyAuth
func CreateNote(c *gin.Context) {
	var note NoteRequest
	if err := c.ShouldBindJSON(&note); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `INSERT INTO note (title, content, created_at, updated_at, user_id, is_public) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := global.Database.Exec(sqlString, note.Title, note.Content, time.Now().Local(),
		time.Now().Local(), c.GetInt("UserId"), c.Query("is_public")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

// UpdateNote godoc
// @Schemes http
// @Description 更新笔记（只有管理员和笔记作者可以更新）
// @Param note body NoteItem true "笔记信息"
// @Param is_public query bool true "是否公开"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /note/update [put]
// @Security ApiKeyAuth
func UpdateNote(c *gin.Context) {
	var note NoteItem
	if err := c.ShouldBindJSON(&note); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT user_id FROM note WHERE id = $1`
	var userId int
	if err := global.Database.Get(&userId, sqlString, note.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if role, _ := c.Get("Role"); userId != c.GetInt("UserId") && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `UPDATE note SET title = $1, content = $2, updated_at = $3, is_public = $4 WHERE id = $5`
	if _, err := global.Database.Exec(sqlString, note.Title, note.Content, time.Now().Local(),
		c.Query("is_public"), note.ID); err != nil {
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
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/like/{id} [post]
// @Security ApiKeyAuth
func LikeNote(c *gin.Context) {
	userId := c.GetInt("UserId")
	noteId := c.Param("id")
	sqlString := `SELECT id FROM note WHERE id = $1`
	var noteIdInt int
	if err := global.Database.Get(&noteIdInt, sqlString, noteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	sqlString = `SELECT id FROM user_like_note WHERE user_id = $1 AND note_id = $2`
	var userLikeNoteId int
	if err := global.Database.Get(&userLikeNoteId, sqlString, userId, noteId); err == nil {
		c.String(http.StatusOK, "点赞成功")
		return
	}
	sqlString = `INSERT INTO user_like_note (user_id, note_id) VALUES ($1, $2)`
	if _, err := global.Database.Exec(sqlString, userId, noteId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
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
	userId := c.GetInt("UserId")
	noteId := c.Param("id")
	// 判断笔记是否存在
	sqlString := `SELECT id FROM note WHERE id = $1`
	var noteIdInt int
	if err := global.Database.Get(&noteIdInt, sqlString, noteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	// 判断是否已经点赞
	sqlString = `SELECT id FROM user_like_note WHERE user_id = $1 AND note_id = $2`
	var userLikeNoteId int
	if err := global.Database.Get(&userLikeNoteId, sqlString, userId, noteId); err != nil {
		c.String(http.StatusOK, "取消点赞成功")
		return
	}
	// 取消点赞
	sqlString = `DELETE FROM user_like_note WHERE user_id = $1 AND note_id = $2`
	if _, err := global.Database.Exec(sqlString, userId, noteId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "取消点赞成功")
}

// FavoriteNote godoc
// @Schemes http
// @Description 收藏笔记
// @Param id path int true "笔记ID"
// @Success 200 {string} string "收藏成功"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/favorite/{id} [post]
// @Security ApiKeyAuth
func FavoriteNote(c *gin.Context) {
	userId := c.GetInt("UserId")
	noteId := c.Param("id")
	sqlString := `SELECT id FROM note WHERE id = $1`
	var noteIdInt int
	if err := global.Database.Get(&noteIdInt, sqlString, noteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	sqlString = `SELECT id FROM user_favorite_note WHERE user_id = $1 AND note_id = $2`
	var userFavoriteNoteId int
	if err := global.Database.Get(&userFavoriteNoteId, sqlString, userId, noteId); err == nil {
		c.String(http.StatusOK, "收藏成功")
		return
	}
	sqlString = `INSERT INTO user_favorite_note (user_id, note_id) VALUES ($1, $2)`
	if _, err := global.Database.Exec(sqlString, userId, noteId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
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
	userId := c.GetInt("UserId")
	noteId := c.Param("id")
	// 判断笔记是否存在
	sqlString := `SELECT id FROM note WHERE id = $1`
	var noteIdInt int
	if err := global.Database.Get(&noteIdInt, sqlString, noteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	// 判断是否已经收藏
	sqlString = `SELECT id FROM user_favorite_note WHERE user_id = $1 AND note_id = $2`
	var userFavoriteNoteId int
	if err := global.Database.Get(&userFavoriteNoteId, sqlString, userId, noteId); err != nil {
		c.String(http.StatusOK, "取消收藏成功")
		return
	}
	// 取消收藏
	sqlString = `DELETE FROM user_favorite_note WHERE user_id = $1 AND note_id = $2`
	if _, err := global.Database.Exec(sqlString, userId, noteId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "取消收藏成功")
}
