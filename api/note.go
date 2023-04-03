package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"net/http"
	"time"
)

type NoteResponse struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
type NoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// GetNotes godoc
// @Schemes http
// @Description 获取当前用户视角下的所有笔记
// @Success 200 {object} []NoteResponse "笔记列表"
// @Failure default {string} string "服务器错误"
// @Router /note/all [get]
func GetNotes(c *gin.Context) {
	var notes []NoteResponse
	var sqlString string
	var err error
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString = `SELECT id, title, content FROM note WHERE is_public = true`
		err = global.Database.Select(&notes, sqlString)
	} else if role == global.USER {
		sqlString = `SELECT id, title, content FROM note WHERE is_public = true OR user_id = $1`
		err = global.Database.Select(&notes, sqlString, c.GetInt("UserId"))
	} else {
		sqlString = `SELECT id, title, content FROM note`
		err = global.Database.Select(&notes, sqlString)
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, notes)
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
	if _, err := global.Database.Exec(sqlString, note.Title, note.Content, time.Now(),
		time.Now(), c.GetInt("UserId"), c.Query("is_public")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

// UpdateNote godoc
// @Schemes http
// @Description 更新笔记（只有管理员和笔记作者可以更新）
// @Param note body NoteResponse true "笔记信息"
// @Param is_public query bool true "是否公开"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /note/update [put]
// @Security ApiKeyAuth
func UpdateNote(c *gin.Context) {
	var note NoteResponse
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
	if _, err := global.Database.Exec(sqlString, note.Title, note.Content, time.Now(),
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
