package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type DiscussionReviewCreateRequest struct {
	Title        string `json:"title"`
	Content      string `json:"content"`
	DiscussionId int    `json:"discussion_id"`
}
type DiscussionReviewResponse struct {
	ID           int              `json:"id"`
	Title        string           `json:"title"`
	Content      string           `json:"content"`
	DiscussionId int              `json:"discussion_id"`
	UserInfo     UserInfoResponse `json:"user_info"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
	IsLiked      bool             `json:"is_liked"`
	LikeCount    int              `json:"like_count"`
}
type AllDiscussionReviewResponse struct {
	TotalCount        int                        `json:"total_count"`
	DiscussionReviews []DiscussionReviewResponse `json:"discussion_reviews"`
}

// AddDiscussionReview godoc
// @Schemes http
// @Description 添加评论
// @Tags DiscussionReview
// @Param review body DiscussionReviewCreateRequest true "评论信息"
// @Success 200 {string} DiscussionReviewResponse "评论信息"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "讨论不存在"
// @Failure default {string} string "服务器错误"
// @Router /discussion_review/add [post]
// @Security ApiKeyAuth
func AddDiscussionReview(c *gin.Context) {
	var request DiscussionReviewCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT * FROM discussion WHERE id = $1`
	var discussion model.Discussion
	if err := global.Database.Get(&discussion, sqlString, request.DiscussionId); err != nil {
		c.String(http.StatusNotFound, "讨论不存在")
		return
	}
	sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
	var count int
	if err := global.Database.Get(&count, sqlString, discussion.GroupId, c.GetInt("user_id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && (count == 0 ||
		discussion.UserId != c.GetInt("user_id") && !discussion.IsPublic) {
		c.String(http.StatusForbidden, "没有权限")
	}
	tx := global.Database.MustBegin()
	sqlString = `INSERT INTO discussion_review (title, content, created_at, updated_at, discussion_id, user_id) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var reviewId int
	if err := global.Database.Get(&reviewId, sqlString, request.Title, request.Content, time.Now().Local(),
		time.Now().Local(), request.DiscussionId, c.GetInt("user_id")); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `SELECT * FROM discussion_review WHERE id = $1`
	var review model.DiscussionReview
	if err := global.Database.Get(&review, sqlString, reviewId); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if err := tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	user := model.User{}
	sqlString = `SELECT name, email, phone, avatar_url, created_at, nick_name FROM "user" WHERE id = $1`
	if err := global.Database.Get(&user, sqlString, discussion.UserId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	userInfo := UserInfoResponse{
		UserId:     user.ID,
		AvatarPath: user.AvatarURL,
		NickName:   user.NickName,
	}
	c.JSON(http.StatusOK, DiscussionReviewResponse{
		ID:           review.ID,
		Title:        review.Title,
		Content:      review.Content,
		DiscussionId: review.DiscussionId,
		UserInfo:     userInfo,
		CreatedAt:    review.CreatedAt,
		UpdatedAt:    review.UpdatedAt,
		IsLiked:      false,
		LikeCount:    0,
	})
}

// RemoveDiscussionReview godoc
// @Schemes http
// @Description 删除评论
// @Tags DiscussionReview
// @Param id path int true "评论ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "评论不存在"
// @Failure default {string} string "服务器错误"
// @Router /discussion_review/remove/{id} [delete]
// @Security ApiKeyAuth
func RemoveDiscussionReview(c *gin.Context) {
	sqlString := `SELECT * FROM discussion_review WHERE id = $1`
	var review model.DiscussionReview
	if err := global.Database.Get(&review, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "评论不存在")
		return
	}
	sqlString = `SELECT * FROM discussion WHERE id = $1`
	var discussion model.Discussion
	if err := global.Database.Get(&discussion, sqlString, review.DiscussionId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `SELECT count(*) FROM "group" WHERE id = $1 AND user_id = $2`
	var count int
	if err := global.Database.Get(&count, sqlString, discussion.GroupId, c.GetInt("user_id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && count == 0 && discussion.UserId != c.GetInt("user_id") {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM discussion_review WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

// GetDiscussionReviews godoc
// @Schemes http
// @Description 获取评论列表
// @Tags DiscussionReview
// @Param discussion_id query int true "讨论ID"
// @Success 200 {string} AllDiscussionReviewResponse "评论列表"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "讨论不存在"
// @Failure default {string} string "服务器错误"
// @Router /discussion_review/get [get]
// @Security ApiKeyAuth
func GetDiscussionReviews(c *gin.Context) {
	var discussion model.Discussion
	discussionId := c.Query("discussion_id")
	sqlString := `SELECT * FROM discussion WHERE id = $1`
	if err := global.Database.Get(&discussion, sqlString, discussionId); err != nil {
		c.String(http.StatusNotFound, "讨论不存在")
		return
	}
	sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
	var count int
	if err := global.Database.Get(&count, sqlString, discussion.GroupId, c.GetInt("user_id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && (count == 0 ||
		discussion.UserId != c.GetInt("user_id") && !discussion.IsPublic) {
		c.String(http.StatusForbidden, "没有权限")
	}
	sqlString = `SELECT * FROM discussion_review WHERE discussion_id = $1`
	var reviews []model.DiscussionReview
	if err := global.Database.Select(&reviews, sqlString, discussionId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var reviewResponses []DiscussionReviewResponse
	for _, review := range reviews {
		user := model.User{}
		sqlString = `SELECT name, email, phone, avatar_url, created_at, nick_name FROM "user" WHERE id = $1`
		if err := global.Database.Get(&user, sqlString, review.UserId); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		userInfo := UserInfoResponse{
			UserId:     user.ID,
			AvatarPath: user.AvatarURL,
			NickName:   user.NickName,
		}
		sqlString = `SELECT count(*) FROM user_like_discussion_review WHERE discussion_review_id = $1 AND user_id = $2`
		var isLiked int
		if err := global.Database.Get(&isLiked, sqlString, review.ID, c.GetInt("user_id")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		reviewResponses = append(reviewResponses, DiscussionReviewResponse{
			ID:           review.ID,
			Title:        review.Title,
			Content:      review.Content,
			DiscussionId: review.DiscussionId,
			UserInfo:     userInfo,
			CreatedAt:    review.CreatedAt,
			UpdatedAt:    review.UpdatedAt,
			IsLiked:      isLiked > 0,
			LikeCount:    review.LikeCount,
		})
	}
	c.JSON(http.StatusOK, AllDiscussionReviewResponse{
		TotalCount:        len(reviewResponses),
		DiscussionReviews: reviewResponses,
	})
}

// LikeDiscussionReview godoc
// @Schemes http
// @Description 点赞评论
// @Tags DiscussionReview
// @Param id path int true "评论ID"
// @Success 200 {string} string "点赞成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "评论不存在"
// @Failure default {string} string "服务器错误"
// @Router /discussion_review/like/{id} [post]
// @Security ApiKeyAuth
func LikeDiscussionReview(c *gin.Context) {
	sqlString := `SELECT * FROM discussion_review WHERE id = $1`
	var review model.DiscussionReview
	if err := global.Database.Get(&review, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "评论不存在")
		return
	}
	sqlString = `SELECT * FROM discussion WHERE id = $1`
	var discussion model.Discussion
	if err := global.Database.Get(&discussion, sqlString, review.DiscussionId); err != nil {
		c.String(http.StatusNotFound, "讨论不存在")
		return
	}
	sqlString = `SELECT count(*) FROM "group" WHERE id = $1 AND user_id = $2`
	var count int
	if err := global.Database.Get(&count, sqlString, discussion.GroupId, c.GetInt("user_id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && (count == 0 ||
		discussion.UserId != c.GetInt("user_id") && !discussion.IsPublic) {
		c.String(http.StatusForbidden, "没有权限")
	}
	sqlString = `SELECT count(*) FROM user_like_discussion_review WHERE user_id = $1 AND discussion_review_id = $2`
	if err := global.Database.Get(&count, sqlString, c.GetInt("user_id"), c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if count > 0 {
		c.String(http.StatusOK, "点赞成功")
		return
	}
	sqlString = `INSERT INTO user_like_discussion_review (user_id, discussion_review_id, created_at) VALUES ($1, $2, now())`
	if _, err := global.Database.Exec(sqlString, c.GetInt("user_id"), c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `UPDATE discussion_review SET like_count = like_count + 1 WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "点赞成功")
}

// UnlikeDiscussionReview godoc
// @Schemes http
// @Description 取消点赞评论
// @Tags DiscussionReview
// @Param id path int true "评论ID"
// @Success 200 {string} string "取消点赞成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "评论不存在"
// @Failure default {string} string "服务器错误"
// @Router /discussion_review/unlike/{id} [post]
// @Security ApiKeyAuth
func UnlikeDiscussionReview(c *gin.Context) {
	sqlString := `SELECT * FROM discussion_review WHERE id = $1`
	var review model.DiscussionReview
	if err := global.Database.Get(&review, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "评论不存在")
		return
	}
	sqlString = `SELECT * FROM discussion WHERE id = $1`
	var discussion model.Discussion
	if err := global.Database.Get(&discussion, sqlString, review.DiscussionId); err != nil {
		c.String(http.StatusNotFound, "讨论不存在")
		return
	}
	sqlString = `DELETE FROM user_like_discussion_review WHERE user_id = $1 AND discussion_review_id = $2`
	if _, err := global.Database.Exec(sqlString, c.GetInt("user_id"), c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `UPDATE discussion_review SET like_count = like_count - 1 WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "取消点赞成功")
}
