package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"net/http"
	"time"
)

type AddReviewRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	NoteId  int    `json:"note_id"`
	UserId  int    `json:"user_id"`
}

type NoteReviewResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	NoteId    int       `json:"note_id"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddNoteReview godoc
// @Schemes http
// @Description 添加评论
// @Param review body AddReviewRequest true "评论信息"
// @Success 200 {string} string "创建成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /note_review/add [post]
// @Security ApiKeyAuth
func AddNoteReview(c *gin.Context) {
	var review AddReviewRequest
	if err := c.ShouldBindJSON(&review); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `INSERT INTO note_review (title, content, note_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := global.Database.Exec(sqlString, review.Title, review.Content, review.NoteId, review.UserId, time.Now(), time.Now())
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "创建成功")
}

// RemoveNoteReview godoc
// @Schemes http
// @Description 删除评论
// @Param review_id query int true "评论id"
// @Success 200 {string} string "删除成功"
// @Failure 404 {string} string "评论不存在"
// @Failure default {string} string "服务器错误"
// @Router /note_review/remove [get]
// @Security ApiKeyAuth
func RemoveNoteReview(c *gin.Context) {
	reviewId := c.Query("review_id")
	// 判断评论是否存在
	var noteId int
	sqlString := `SELECT note_id FROM note_review WHERE id = $1`
	err := global.Database.QueryRow(sqlString, reviewId).Scan(&noteId)
	if err != nil {
		c.String(http.StatusNotFound, "评论不存在")
		return
	}

	// 判断是否有权限
	// 1. 笔记作者是当前用户
	UserId := c.GetInt("UserId")
	sqlString = `SELECT note_id FROM note_review WHERE id = $1`
	err = global.Database.QueryRow(sqlString, reviewId).Scan(&noteId)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `SELECT user_id FROM note WHERE id = $1`
	err = global.Database.QueryRow(sqlString, noteId).Scan(&noteId)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if noteId != UserId {
		c.String(http.StatusForbidden, "无权限")
		return
	}
	// 2. 评论作者是当前用户
	sqlString = `SELECT user_id FROM note_review WHERE id = $1`
	err = global.Database.QueryRow(sqlString, reviewId).Scan(&noteId)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if noteId != UserId {
		c.String(http.StatusForbidden, "无权限")
		return
	}

	// 删除评论
	sqlString = `DELETE FROM note_review WHERE id = $1`
	_, err = global.Database.Exec(sqlString, reviewId)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

// GetNoteReviews godoc
// @Schemes http
// @Description 获取笔记的评论
// @Param note_id query int true "笔记id"
// @Success 200 {object} []NoteReviewResponse "评论列表"
// @Failure default {string} string "服务器错误"
// @Router /note_review/get [get]
func GetNoteReviews(c *gin.Context) {
	var reviews []NoteReviewResponse
	noteId := c.Query("note_id")
	sqlString := `SELECT id, title, content, note_id, user_id, created_at, updated_at FROM note_review WHERE note_id = $1`
	err := global.Database.Select(&reviews, sqlString, noteId)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, reviews)
}
