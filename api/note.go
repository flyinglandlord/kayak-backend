package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"strconv"
	"time"
)

type NoteFilter struct {
	ID         *int  `json:"id" form:"id"`
	UserId     *int  `json:"user_id" form:"user_id"`
	IsLiked    *bool `json:"is_liked" form:"is_liked"`
	IsFavorite *bool `json:"is_favorite" form:"is_favorite"`
	Offset     *int  `json:"offset" form:"offset"`
	Limit      *int  `json:"limit" form:"limit"`
	SortByLike *bool `json:"sort_by_like" form:"sort_by_like"`
}
type NoteResponse struct {
	ID            int       `json:"id" db:"id"`
	UserId        int       `json:"user_id" db:"user_id"`
	Title         string    `json:"title" db:"title"`
	Content       string    `json:"content" db:"content"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	IsLiked       bool      `json:"is_liked"`
	LikeCount     int       `json:"like_count"`
	IsFavorite    bool      `json:"is_favorite"`
	FavoriteCount int       `json:"favorite_count"`
	IsPublic      bool      `json:"is_public" db:"is_public"`
}
type AllNoteResponse struct {
	TotalCount int            `json:"total_count"`
	Notes      []NoteResponse `json:"notes"`
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
// @Description 获取符合filter要求的当前用户视角下的所有笔记
// @Tags note
// @Param filter query NoteFilter false "筛选条件"
// @Success 200 {object} AllNoteResponse "笔记列表"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /note/all [get]
// @Security ApiKeyAuth
func GetNotes(c *gin.Context) {
	sqlString := `SELECT note.*, count(id) like_count FROM note, user_like_note WHERE note.id = user_like_note.note_id group by id`
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString += ` WHERE is_public = true`
	} else if role == global.USER {
		sqlString += ` WHERE (is_public = true OR user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
	} else {
		sqlString += ` WHERE 1 = 1`
	}
	var filter NoteFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	if filter.ID != nil {
		sqlString += ` AND id = ` + strconv.Itoa(*filter.ID)
	}
	if filter.IsLiked != nil {
		if *filter.IsLiked {
			sqlString += ` AND id IN (SELECT note_id FROM user_like_note WHERE user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
		} else {
			sqlString += ` AND id NOT IN (SELECT note_id FROM user_like_note WHERE user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
		}
	}
	if filter.IsFavorite != nil {
		if *filter.IsFavorite {
			sqlString += ` AND id IN (SELECT note_id FROM user_favorite_note WHERE user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
		} else {
			sqlString += ` AND id NOT IN (SELECT note_id FROM user_favorite_note WHERE user_id = ` + strconv.Itoa(c.GetInt("UserId")) + `)`
		}
	}
	if filter.UserId != nil {
		sqlString += fmt.Sprintf(` AND user_id = %d`, *filter.UserId)
	}
	if filter.SortByLike != nil {
		if *filter.SortByLike {
			sqlString += ` ORDER BY like_count DESC`
		} else {
			sqlString += ` ORDER BY created_at DESC`
		}
	}
	if filter.Limit != nil {
		sqlString += ` LIMIT ` + strconv.Itoa(*filter.Limit)
	}
	if filter.Offset != nil {
		sqlString += ` OFFSET ` + strconv.Itoa(*filter.Offset)
	}
	var notes []model.Note
	if err := global.Database.Select(&notes, sqlString); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var noteResponses []NoteResponse
	for _, note := range notes {
		var isLiked, isFavorite int
		var likeCount, favoriteCount int
		sqlString = `SELECT COUNT(*) FROM user_like_note WHERE note_id = $1 AND user_id = $2`
		if err := global.Database.Get(&isLiked, sqlString, note.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		sqlString = `SELECT COUNT(*) FROM user_like_note WHERE note_id = $1`
		if err := global.Database.Get(&likeCount, sqlString, note.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		sqlString = `SELECT COUNT(*) FROM user_favorite_note WHERE note_id = $1 AND user_id = $2`
		if err := global.Database.Get(&isFavorite, sqlString, note.ID, c.GetInt("UserId")); err != nil {
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
			UserId:        note.UserId,
			Title:         note.Title,
			Content:       note.Content,
			CreatedAt:     note.CreatedAt,
			IsLiked:       isLiked > 0,
			LikeCount:     likeCount,
			IsFavorite:    isFavorite > 0,
			FavoriteCount: favoriteCount,
			IsPublic:      note.IsPublic,
		})
	}
	c.JSON(http.StatusOK, AllNoteResponse{
		TotalCount: len(noteResponses),
		Notes:      noteResponses,
	})
}

// CreateNote godoc
// @Schemes http
// @Description 创建笔记
// @Tags note
// @Param note body NoteCreateRequest true "笔记信息"
// @Success 200 {object} NoteResponse "笔记信息"
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
	tx := global.Database.MustBegin()
	sqlString := `INSERT INTO note (title, content, created_at, updated_at, user_id, is_public) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var noteId int
	if err := global.Database.Get(&noteId, sqlString, request.Title, request.Content, time.Now().Local(),
		time.Now().Local(), c.GetInt("UserId"), request.IsPublic); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	sqlString = `SELECT * FROM note WHERE id = $1`
	var note model.Note
	if err := global.Database.Get(&note, sqlString, noteId); err != nil {
		_ = tx.Rollback()
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if err := tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, NoteResponse{
		ID:            note.ID,
		UserId:        note.UserId,
		Title:         note.Title,
		Content:       note.Content,
		CreatedAt:     note.CreatedAt,
		IsLiked:       false,
		LikeCount:     0,
		IsFavorite:    false,
		FavoriteCount: 0,
		IsPublic:      note.IsPublic,
	})
}

// UpdateNote godoc
// @Schemes http
// @Description 更新笔记（只有管理员和笔记作者可以更新）(可以只传需要更新的字段)
// @Tags note
// @Param note body NoteUpdateRequest true "笔记信息"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/update [put]
// @Security ApiKeyAuth
func UpdateNote(c *gin.Context) {
	var request NoteUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "请求解析失败")
		return
	}
	sqlString := `SELECT * FROM note WHERE id = $1`
	var note model.Note
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
// @Tags note
// @Param id path int true "笔记ID"
// @Success 200 {string} string "删除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/delete/{id} [delete]
// @Security ApiKeyAuth
func DeleteNote(c *gin.Context) {
	sqlString := `SELECT user_id FROM note WHERE id = $1`
	var noteUserId int
	if err := global.Database.Get(&noteUserId, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); c.GetInt("UserId") != noteUserId && role != global.ADMIN {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `DELETE FROM note WHERE id = $1`
	if _, err := global.Database.Exec(sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "删除成功")
}

// LikeNote godoc
// @Schemes http
// @Description 点赞笔记
// @Tags note
// @Param id path int true "笔记ID"
// @Success 200 {string} string "点赞成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/like/{id} [post]
// @Security ApiKeyAuth
func LikeNote(c *gin.Context) {
	sqlString := `SELECT * FROM note WHERE id = $1`
	var note model.Note
	if err := global.Database.Get(&note, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && note.UserId != c.GetInt("UserId") && !note.IsPublic {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `INSERT INTO user_like_note (user_id, note_id, created_at) VALUES ($1, $2, $3) ON CONFLICT (user_id, note_id) do update set created_at = $3`
	if _, err := global.Database.Exec(sqlString, c.GetInt("UserId"), c.Param("id"), time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "点赞成功")
}

// UnlikeNote godoc
// @Schemes http
// @Description 取消点赞笔记
// @Tags note
// @Param id path int true "笔记ID"
// @Success 200 {string} string "取消点赞成功"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/unlike/{id} [post]
// @Security ApiKeyAuth
func UnlikeNote(c *gin.Context) {
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

// AddProblemToNote godoc
// @Schemes http
// @Description 将题目添加到笔记
// @Tags note
// @Param id path int true "笔记ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "添加成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在或题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/add_problem/{id} [post]
// @Security ApiKeyAuth
func AddProblemToNote(c *gin.Context) {
	sqlString := `SELECT * FROM note WHERE id = $1`
	var note model.Note
	if err := global.Database.Get(&note, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && note.UserId != c.GetInt("UserId") {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT * FROM problem_type WHERE id = $1`
	var problem model.ProblemType
	if err := global.Database.Get(&problem, sqlString, c.Query("problem_id")); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	sqlString = `INSERT INTO note_problem (note_id, problem_id, created_at) VALUES ($1, $2, $3) ON CONFLICT (note_id, problem_id) do update set created_at = $3`
	if _, err := global.Database.Exec(sqlString, c.Param("id"), c.Query("problem_id"), time.Now().Local()); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "添加成功")
}

// RemoveProblemFromNote godoc
// @Schemes http
// @Description 将题目从笔记中移除
// @Tags note
// @Param id path int true "笔记ID"
// @Param problem_id query int true "题目ID"
// @Success 200 {string} string "移除成功"
// @Failure 403 {string} string "没有权限"
// @Failure 404 {string} string "笔记不存在或题目不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/remove_problem/{id} [delete]
// @Security ApiKeyAuth
func RemoveProblemFromNote(c *gin.Context) {
	sqlString := `SELECT * FROM note WHERE id = $1`
	var note model.Note
	if err := global.Database.Get(&note, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	if role, _ := c.Get("Role"); role != global.ADMIN && note.UserId != c.GetInt("UserId") {
		c.String(http.StatusForbidden, "没有权限")
		return
	}
	sqlString = `SELECT * FROM problem_type WHERE id = $1`
	var problem model.ProblemType
	if err := global.Database.Get(&problem, sqlString, c.Query("problem_id")); err != nil {
		c.String(http.StatusNotFound, "题目不存在")
		return
	}
	sqlString = `DELETE FROM note_problem WHERE note_id = $1 AND problem_id = $2`
	if _, err := global.Database.Exec(sqlString, c.Param("id"), c.Query("problem_id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.String(http.StatusOK, "移除成功")
}

// GetNoteProblemList godoc
// @Schemes http
// @Description 获取笔记中的题目列表
// @Tags note
// @Param id path int true "笔记ID"
// @Param offset query int false "页数"
// @Param limit query int false "每页数量"
// @Success 200 {string} string "获取成功"
// @Failure 404 {string} string "笔记不存在"
// @Failure default {string} string "服务器错误"
// @Router /note/problem_list/{id} [get]
// @Security ApiKeyAuth
func GetNoteProblemList(c *gin.Context) {
	sqlString := `SELECT * FROM note WHERE id = $1`
	var note model.Note
	if err := global.Database.Get(&note, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "笔记不存在")
		return
	}
	sqlString = `SELECT * FROM problem_type WHERE id IN (SELECT problem_id FROM note_problem WHERE note_id = $1) ORDER BY id DESC`
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	var problems []model.ProblemType
	if err := global.Database.Select(&problems, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	if len(problems) == 0 {
		c.JSON(http.StatusOK, make([]model.ProblemType, 0))
		return
	}
	c.JSON(http.StatusOK, problems)
}
