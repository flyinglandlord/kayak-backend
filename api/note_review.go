package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type NoteReviewCreateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	NoteId  int    `json:"note_id"`
}
type NoteReviewResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	NoteId    int       `json:"note_id"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsLiked   bool      `json:"is_liked"`
	LikeCount int       `json:"like_count"`
}
type AllNoteReviewResponse struct {
	TotalCount  int                  `json:"total_count"`
	NoteReviews []NoteReviewResponse `json:"note_reviews"`
}

// AddNoteReview godoc
// @Schemes http
// @Description 添加评论
// @Param review body NoteReviewCreateRequest true "评论信息"
// @Success 200 {string} NoteReviewResponse "评论信息"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note_review/add [post]
// @Security ApiKeyAuth
func AddNoteReview(c *gin.Context) {
	var request NoteReviewCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT * FROM note WHERE id = $1`
	var note model.Note
	if err := global.Database.Get(&note, sqlString, request.NoteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && note.UserId != c.GetInt("UserId") && !note.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	tx := global.Database.MustBegin()
	sqlString = `INSERT INTO note_review (title, content, note_id, user_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var reviewId int
	if err := global.Database.Get(&reviewId, sqlString, request.Title, request.Content, request.NoteId,
		c.GetInt("UserId"), time.Now().Local(), time.Now().Local()); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `SELECT * FROM note_review WHERE id = $1`
	var review model.NoteReview
	if err := global.Database.Get(&review, sqlString, reviewId); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if err := tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, NoteReviewResponse{
		ID:        review.ID,
		Title:     review.Title,
		Content:   review.Content,
		NoteId:    review.NoteId,
		UserId:    review.UserId,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		IsLiked:   false,
		LikeCount: 0,
	})
}

// RemoveNoteReview godoc
// @Schemes http
// @Description 删除评论
// @Param id path int true "评论id"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "评论不存在"/"笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note_review/remove/{id} [delete]
// @Security ApiKeyAuth
func RemoveNoteReview(c *gin.Context) {
	sqlString := `SELECT * FROM note_review WHERE id = $1`
	var noteReview model.NoteReview
	if err := global.Database.Get(&noteReview, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "评论不存在")
		return
	}
	sqlString = `SELECT * FROM note WHERE id = $1`
	var note model.Note
	if err := global.Database.Get(&note, sqlString, noteReview.NoteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && noteReview.UserId != c.GetInt("UserId") &&
		note.UserId != c.GetInt("UserId") {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM note_review WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

// GetNoteReviews godoc
// @Schemes http
// @Description 获取笔记的评论
// @Param note_id query int true "笔记id"
// @Success 200 {object} AllNoteReviewResponse "评论列表"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note_review/get [get]
// @Security ApiKeyAuth
func GetNoteReviews(c *gin.Context) {
	var note model.Note
	noteId := c.Query("note_id")
	sqlString := `SELECT * FROM note WHERE id = $1`
	if err := global.Database.Get(&note, sqlString, noteId); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && note.UserId != c.GetInt("UserId") && !note.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT * FROM note_review WHERE note_id = $1`
	var reviews []model.NoteReview
	if err := global.Database.Select(&reviews, sqlString, noteId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var reviewResponses []NoteReviewResponse
	for _, review := range reviews {
		var isLiked int
		sqlString = `SELECT COUNT(*) FROM user_like_note_review WHERE note_review_id = $1 AND user_id = $2`
		if err := global.Database.Get(&isLiked, sqlString, review.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		var likeCount int
		sqlString = `SELECT COUNT(*) FROM user_like_note_review WHERE note_review_id = $1`
		if err := global.Database.Get(&likeCount, sqlString, review.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		reviewResponses = append(reviewResponses, NoteReviewResponse{
			ID:        review.ID,
			Title:     review.Title,
			Content:   review.Content,
			NoteId:    review.NoteId,
			UserId:    review.UserId,
			CreatedAt: review.CreatedAt,
			UpdatedAt: review.UpdatedAt,
			IsLiked:   isLiked > 0,
			LikeCount: likeCount,
		})
	}
	c.JSON(http.StatusOK, AllNoteReviewResponse{
		TotalCount:  len(reviewResponses),
		NoteReviews: reviewResponses,
	})
}
