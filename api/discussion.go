package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"time"
)

type DiscussionFilter struct {
	GroupId    int   `json:"group_id" form:"group_id"`
	ID         *int  `json:"id" form:"id"`
	UserId     *int  `json:"user_id" form:"user_id"`
	IsLiked    *bool `json:"is_liked" form:"is_liked"`
	Offset     *int  `json:"offset" form:"offset"`
	Limit      *int  `json:"limit" form:"limit"`
	SortByLike *bool `json:"sort_by_like" form:"sort_by_like"`
}
type DiscussionResponse struct {
	ID        int              `json:"id" db:"id"`
	Title     string           `json:"title" db:"title"`
	Content   string           `json:"content" db:"content"`
	UserInfo  UserInfoResponse `json:"user_info" db:"user_info"`
	GroupId   int              `json:"group_id" db:"group_id"`
	CreatedAt string           `json:"created_at" db:"created_at"`
	UpdatedAt string           `json:"updated_at" db:"updated_at"`
	IsPublic  bool             `json:"is_public" db:"is_public"`
	IsLiked   bool             `json:"is_liked" db:"is_liked"`
	LikeCount int              `json:"like_count" db:"like_count"`
}
type AllDiscussionResponse struct {
	TotalCount  int                  `json:"total_count"`
	Discussions []DiscussionResponse `json:"discussions"`
}
type DiscussionCreateRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	GroupId  int    `json:"group_id"`
	IsPublic bool   `json:"is_public"`
}
type DiscussionUpdateRequest struct {
	ID       int     `json:"id"`
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	IsPublic *bool   `json:"is_public"`
}

// GetDiscussions godoc
// @Schemes http
// @Description 获取符合filter要求的当前用户视角下的所有讨论
// @Tags Discussion
// @Param filter query DiscussionFilter false "筛选条件"
// @Success 200 {object} AllDiscussionResponse "讨论列表"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /discussion/all [get]
// @Security ApiKeyAuth
func GetDiscussions(c *gin.Context) {
	var filter DiscussionFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, "服务器错误")
		return
	}
	sqlString := `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
	var count int
	if err := global.Database.Get(&count, sqlString, filter.GroupId, c.GetInt("user_id")); err != nil {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT * FROM discussion WHERE group_id = $1 AND (is_public = true OR user_id = $2)`
	if filter.ID != nil {
		sqlString += fmt.Sprint(" AND id = ", *filter.ID)
	}
	if filter.UserId != nil {
		sqlString += fmt.Sprint(" AND user_id = ", *filter.UserId)
	}
	if filter.IsLiked != nil {
		if *filter.IsLiked {
			sqlString += fmt.Sprint(" AND id IN (SELECT discussion_id FROM user_like_discussion WHERE user_id = ", c.GetInt("user_id"), ")")
		} else {
			sqlString += fmt.Sprint(" AND id NOT IN (SELECT discussion_id FROM user_like_discussion WHERE user_id = ", c.GetInt("user_id"), ")")
		}
	}
	if filter.SortByLike != nil {
		if *filter.SortByLike {
			sqlString += " ORDER BY like_count DESC"
		} else {
			sqlString += " ORDER BY created_at DESC"
		}
	}
	if filter.Limit != nil {
		sqlString += fmt.Sprint(" LIMIT ", *filter.Limit)
	}
	if filter.Offset != nil {
		sqlString += fmt.Sprint(" OFFSET ", *filter.Offset)
	}
	var discussions []model.Discussion
	if err := global.Database.Select(&discussions, sqlString, filter.GroupId, c.GetInt("UserId")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var discussionResponses []DiscussionResponse
	for _, discussion := range discussions {
		var isLiked int
		sqlString = `SELECT count(*) FROM user_like_discussion WHERE user_id = $1 AND discussion_id = $2`
		if err := global.Database.Get(&isLiked, sqlString, c.GetInt("user_id"), discussion.ID); err != nil {
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
		discussionResponses = append(discussionResponses, DiscussionResponse{
			ID:        discussion.ID,
			Title:     discussion.Title,
			Content:   discussion.Content,
			UserInfo:  userInfo,
			GroupId:   discussion.GroupId,
			CreatedAt: discussion.CreatedAt,
			UpdatedAt: discussion.UpdatedAt,
			IsPublic:  discussion.IsPublic,
			IsLiked:   isLiked > 0,
			LikeCount: discussion.LikeCount,
		})
	}
	c.JSON(http.StatusOK, AllDiscussionResponse{
		TotalCount:  len(discussionResponses),
		Discussions: discussionResponses,
	})
}

// CreateDiscussion godoc
// @Schemes http
// @Tags Discussion
// @Param note body DiscussionCreateRequest true "讨论信息"
// @Success 200 {object} DiscussionResponse "笔记信息"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure default {string} string "服务器错误"
// @Router /discussion/create [post]
// @Security ApiKeyAuth
func CreateDiscussion(c *gin.Context) {
	var request DiscussionCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "服务器错误")
		return
	}
	sqlString := `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
	var count int
	if err := global.Database.Get(&count, sqlString, request.GroupId, c.GetInt("user_id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if count == 0 {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO discussion (title, content, user_id, group_id, created_at, updated_at, is_public) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var discussionId int
	if err := global.Database.Get(&discussionId, sqlString, request.Title, request.Content, c.GetInt("user_id"),
		request.GroupId, time.Now().Local(), time.Now().Local(), request.IsPublic); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var discussion model.Discussion
	sqlString = `SELECT * FROM discussion WHERE id = $1`
	if err := global.Database.Get(&discussion, sqlString, discussionId); err != nil {
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
	c.JSON(http.StatusOK, DiscussionResponse{
		ID:        discussion.ID,
		Title:     discussion.Title,
		Content:   discussion.Content,
		UserInfo:  userInfo,
		GroupId:   discussion.GroupId,
		CreatedAt: discussion.CreatedAt,
		UpdatedAt: discussion.UpdatedAt,
		IsPublic:  discussion.IsPublic,
		IsLiked:   false,
		LikeCount: discussion.LikeCount,
	})
}

// UpdateDiscussion godoc
// @Schemes http
// @Description 更新讨论（只有创建者可以修改）
// @Tags Discussion
// @Param note body NoteUpdateRequest true "讨论信息"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "讨论不存在"
// @Failure default {string} string "服务器错误"
// @Router /discussion/update [put]
// @Security ApiKeyAuth
func UpdateDiscussion(c *gin.Context) {
	var request DiscussionUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT * FROM discussion WHERE id = $1`
	var discussion model.Discussion
	if err := global.Database.Get(&discussion, sqlString, request.ID); err != nil {
		c.String(http.StatusNotFound, "讨论不存在")
		return
	}
	if role, _ := c.Get("Role"); discussion.UserId != c.GetInt("UserId") && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	if request.Title == nil {
		request.Title = &discussion.Title
	}
	if request.Content == nil {
		request.Content = &discussion.Content
	}
	if request.IsPublic == nil {
		request.IsPublic = &discussion.IsPublic
	}
	sqlString = `UPDATE discussion SET title = $1, content = $2, updated_at = $3, is_public = $4 WHERE id = $5`
	if _, err := global.Database.Exec(sqlString, request.Title, request.Content,
		time.Now().Local(), request.IsPublic, request.ID); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "更新成功")
}

// DeleteDiscussion godoc
// @Schemes http
// @Description 删除讨论（只有创建者可以删除）
// @Tags Discussion
// @Param id path int true "讨论ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "讨论不存在"
// @Failure default {string} string "服务器错误"
// @Router /discussion/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteDiscussion(c *gin.Context) {
	sqlString := `SELECT user_id FROM discussion WHERE id = $1`
	var discussionUserId int
	if err := global.Database.Get(&discussionUserId, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "讨论不存在")
		return
	}
	if role, _ := c.Get("Role"); c.GetInt("UserId") != discussionUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM discussion WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

// LikeDiscussion godoc
// @Schemes http
// @Description 点赞讨论
// @Tags Discussion
// @Param id path int true "讨论ID"
// @Success 200 {string} string "点赞成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "讨论不存在"
// @Failure default {string} string "服务器错误"
// @Router /discussion/like/{id} [post]
// @Security ApiKeyAuth
func LikeDiscussion(c *gin.Context) {
	sqlString := `SELECT * FROM discussion WHERE id = $1`
	var discussion model.Discussion
	if err := global.Database.Get(&discussion, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "讨论不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && discussion.UserId != c.GetInt("UserId") && !discussion.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	var count int
	sqlString = `SELECT count(*) FROM group_member WHERE user_id = $1 AND group_id = $2`
	if err := global.Database.Get(&count, sqlString, c.GetInt("UserId"), discussion.GroupId); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if count == 0 {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO user_like_discussion (user_id, discussion_id, created_at) VALUES ($1, $2, $3) ON CONFLICT (user_id, discussion_id) do update set created_at = $3`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Param("id"), time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "点赞成功")
}

// UnlikeDiscussion godoc
// @Schemes http
// @Description 取消点赞讨论
// @Tags Discussion
// @Param id path int true "讨论ID"
// @Success 200 {string} string "取消点赞成功"
// @Failure 404 {string} string "讨论不存在"
// @Failure default {string} string "服务器错误"
// @Router /discussion/unlike/{id} [post]
// @Security ApiKeyAuth
func UnlikeDiscussion(c *gin.Context) {
	sqlString := `SELECT * FROM note WHERE id = $1`
	var note model.Note
	if err := global.Database.Get(&note, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	sqlString = `DELETE FROM user_like_note WHERE user_id = $1 AND note_id = $2`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "取消点赞成功")
}
