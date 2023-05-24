package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
)

// This file will implement some special APIs, which is not in our expectation.
// I do not know why we need it, but front-end requires it.

type AllWrongProblemSet struct {
	TotalCount  int `json:"total_count"`
	ProblemSets []struct {
		Name         string `json:"name"`
		ID           int    `json:"id"`
		ProblemCount int    `json:"problem_count"`
	}
}

type AllFavoriteProblemSet struct {
	TotalCount  int `json:"total_count"`
	ProblemSets []struct {
		Name         string `json:"name"`
		ID           int    `json:"id"`
		ProblemCount int    `json:"problem_count"`
	}
}

// GetWrongProblemSet godoc
// @Schemes http
// @Description 获取错题集
// @Tags Special
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} AllWrongProblemSet "错题集"
// @Failure default {string} string "服务器错误"
// @Router /special/wrong_problem_set [get]
// @Security ApiKeyAuth
func GetWrongProblemSet(c *gin.Context) {
	sqlString := `SELECT * FROM problem_set`
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString += ` WHERE is_public = true`
	} else if role == global.USER {
		sqlString += ` WHERE (is_public = true OR user_id = ` + fmt.Sprintf("%d", c.GetInt("UserId")) + `)`
	} else {
		sqlString += ` WHERE 1 = 1`
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	var problemSets []model.ProblemSet
	if err := global.Database.Select(&problemSets, sqlString); err != nil {
		c.String(http.StatusBadRequest, "服务器错误")
		return
	}
	var allWrongProblemSet AllWrongProblemSet
	for _, problemSet := range problemSets {
		var problemCount int
		sqlString = `SELECT COUNT(*) FROM user_wrong_record WHERE user_id = $1 AND problem_id IN (SELECT problem_id FROM problem_in_problem_set WHERE problem_set_id = $2)`
		if err := global.Database.Get(&problemCount, sqlString, c.GetInt("UserId"), problemSet.ID); err != nil {
			c.String(http.StatusBadRequest, "服务器错误")
			return
		}
		if problemCount > 0 {
			allWrongProblemSet.ProblemSets = append(allWrongProblemSet.ProblemSets, struct {
				Name         string `json:"name"`
				ID           int    `json:"id"`
				ProblemCount int    `json:"problem_count"`
			}{Name: problemSet.Name, ID: problemSet.ID, ProblemCount: problemCount})
		}
	}
	allWrongProblemSet.TotalCount = len(allWrongProblemSet.ProblemSets)
	c.JSON(http.StatusOK, allWrongProblemSet)
}

// GetFavoriteProblemSet godoc
// @Schemes http
// @Description 获取含有收藏题目的题集
// @Tags Special
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"
// @Success 200 {object} AllFavoriteProblemSet "收藏题集"
// @Failure default {string} string "服务器错误"
// @Router /special/favorite_problem_set [get]
// @Security ApiKeyAuth
func GetFavoriteProblemSet(c *gin.Context) {
	sqlString := `SELECT * FROM problem_set`
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString += ` WHERE is_public = true`
	} else if role == global.USER {
		sqlString += ` WHERE (is_public = true OR user_id = ` + fmt.Sprintf("%d", c.GetInt("UserId")) + `)`
	} else {
		sqlString += ` WHERE 1 = 1`
	}
	if c.Query("limit") != "" {
		sqlString += ` LIMIT ` + c.Query("limit")
	}
	if c.Query("offset") != "" {
		sqlString += ` OFFSET ` + c.Query("offset")
	}
	var problemSets []model.ProblemSet
	if err := global.Database.Select(&problemSets, sqlString); err != nil {
		c.String(http.StatusBadRequest, "服务器错误")
		return
	}
	var allFavoriteProblemSet AllFavoriteProblemSet
	for _, problemSet := range problemSets {
		var problemCount int
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem WHERE user_id = $1 AND problem_id IN (SELECT problem_id FROM problem_in_problem_set WHERE problem_set_id = $2)`
		if err := global.Database.Get(&problemCount, sqlString, c.GetInt("UserId"), problemSet.ID); err != nil {
			c.String(http.StatusBadRequest, "服务器错误")
			return
		}
		if problemCount > 0 {
			allFavoriteProblemSet.ProblemSets = append(allFavoriteProblemSet.ProblemSets, struct {
				Name         string `json:"name"`
				ID           int    `json:"id"`
				ProblemCount int    `json:"problem_count"`
			}{Name: problemSet.Name, ID: problemSet.ID, ProblemCount: problemCount})
		}
	}
	allFavoriteProblemSet.TotalCount = len(allFavoriteProblemSet.ProblemSets)
	c.JSON(http.StatusOK, allFavoriteProblemSet)
}

// GetFeaturedProblemSet godoc
// @Schemes http
// @Description 获取精选题集
// @Tags Special
// @Success 200 {object} AllProblemSetResponse "题集列表"
// @Failure default {string} string "服务器错误"
// @Router /special/featured_problem_set [get]
// @Security ApiKeyAuth
func GetFeaturedProblemSet(c *gin.Context) {
	sqlString := `SELECT ps.id, ps.name, ps.description, ps.created_at, ps.updated_at, ps.user_id, ps.is_public, ps.group_id
		FROM problem_set ps LEFT JOIN user_favorite_problem_set ufps ON ps.id = ufps.problem_set_id WHERE ps.is_public = true GROUP BY ps.id ORDER BY count(*) DESC LIMIT 6`
	var problemSets []model.ProblemSet
	if err := global.Database.Select(&problemSets, sqlString); err != nil {
		c.String(http.StatusBadRequest, "服务器错误")
		return
	}
	var problemSetResponses []ProblemSetResponse
	for _, problemSet := range problemSets {
		var problemCount int
		sqlString = `SELECT COUNT(*) FROM problem_in_problem_set WHERE problem_set_id = $1`
		if err := global.Database.Get(&problemCount, sqlString, problemSet.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem_set WHERE problem_set_id = $1 AND user_id = $2`
		var isFavorite int
		if err := global.Database.Get(&isFavorite, sqlString, problemSet.ID, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		var favoriteCount int
		sqlString = `SELECT COUNT(*) FROM user_favorite_problem_set WHERE problem_set_id = $1`
		if err := global.Database.Get(&favoriteCount, sqlString, problemSet.ID); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		user := model.User{}
		sqlString = `SELECT id, avatar_url, nick_name FROM "user" WHERE id = $1`
		if err := global.Database.Get(&user, sqlString, problemSet.UserId); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		userInfo := UserInfoResponse{
			UserId:     user.ID,
			AvatarPath: user.AvatarURL,
			NickName:   user.NickName,
		}
		problemSetResponses = append(problemSetResponses, ProblemSetResponse{
			ID:            problemSet.ID,
			Name:          problemSet.Name,
			Description:   problemSet.Description,
			CreatedAt:     problemSet.CreatedAt,
			UpdatedAt:     problemSet.UpdatedAt,
			ProblemCount:  problemCount,
			IsFavorite:    isFavorite > 0,
			FavoriteCount: favoriteCount,
			UserId:        problemSet.UserId,
			UserInfo:      userInfo,
			IsPublic:      problemSet.IsPublic,
			GroupId:       problemSet.GroupId,
			AreaId:        problemSet.AreaId,
		})
	}
	c.JSON(http.StatusOK, AllProblemSetResponse{
		TotalCount: len(problemSetResponses),
		ProblemSet: problemSetResponses,
	})
}

// GetFeaturedNote godoc
// @Schemes http
// @Description 获取精选笔记
// @Tags Special
// @Success 200 {object} AllNoteResponse "笔记列表"
// @Failure default {string} string "服务器错误"
// @Router /special/featured_note [get]
// @Security ApiKeyAuth
func GetFeaturedNote(c *gin.Context) {
	sqlString := `SELECT n.id, n.title, n.content, n.created_at, n.updated_at, n.user_id, n.is_public 
		FROM note n LEFT JOIN user_favorite_note ufn ON n.id = ufn.note_id WHERE n.is_public = true GROUP BY n.id ORDER BY count(*) DESC LIMIT 6`
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

// GetFeaturedGroup godoc
// @Schemes http
// @Description 获取精选小组
// @Tags Special
// @Success 200 {object} AllGroupResponse "小组列表"
// @Failure default {string} string "服务器错误"
// @Router /special/featured_group [get]
// @Security ApiKeyAuth
func GetFeaturedGroup(c *gin.Context) {
	sqlString := `SELECT g.id, g.name, g.description, g.created_at, g.user_id FROM "group" g LEFT JOIN group_member gm 
    	ON g.id = gm.group_id GROUP BY g.id ORDER BY count(*) DESC LIMIT 6`
	var groups []model.Group
	if err := global.Database.Select(&groups, sqlString); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	var groupResponses []GroupResponse
	for _, group := range groups {
		user := model.User{}
		sqlString = `SELECT id, avatar_url, nick_name FROM "user" WHERE id = $1`
		if err := global.Database.Get(&user, sqlString, group.UserId); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		userInfo := UserInfoResponse{
			UserId:     user.ID,
			AvatarPath: user.AvatarURL,
			NickName:   user.NickName,
		}
		var count int
		sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1`
		if err := global.Database.Get(&count, sqlString, group.Id); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		groupResponses = append(groupResponses, GroupResponse{
			Id:          group.Id,
			Name:        group.Name,
			Description: group.Description,
			UserId:      group.UserId,
			UserInfo:    userInfo,
			MemberCount: count,
			CreatedAt:   group.CreatedAt,
			AreaId:      group.AreaId,
		})
	}
	c.JSON(http.StatusOK, AllGroupResponse{
		TotalCount: len(groupResponses),
		Group:      groupResponses,
	})
}
