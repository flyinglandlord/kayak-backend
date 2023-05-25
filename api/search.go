package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
)

type SearchRequest struct {
	Keyword string `json:"keyword" form:"keyword"`
	Limit   int    `json:"limit" form:"limit"`
	Offset  int    `json:"offset" form:"offset"`
}

// SearchProblemSets godoc
// @Schemes http
// @Description 搜索题集
// @Tags Search
// @Param search query SearchRequest true "搜索信息"
// @Success 200 {object} AllProblemSetResponse "题集列表"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /search/problem_set [get]
// @Security ApiKeyAuth
func SearchProblemSets(c *gin.Context) {
	var request SearchRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.String(http.StatusBadRequest, "请求错误")
		return
	}
	sqlString := fmt.Sprintf(`SELECT
			problem_set.*
		FROM
			problem_set, 
			to_tsvector(problem_set.name || problem_set.description) document,
			to_tsquery($1) query,
			NULLIF(ts_rank(to_tsvector(problem_set.name), query), 0) rank_name,
			NULLIF(ts_rank(to_tsvector(problem_set.description), query), 0) rank_description,
			SIMILARITY($1, problem_set.name || problem_set.description) similarity`)
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString += ` WHERE is_public = true`
	} else if role == global.USER {
		sqlString += fmt.Sprint(` WHERE (is_public = true OR (user_id = `, c.GetInt("UserId"), ` AND group_id = 0) 
			OR (`, c.GetInt("UserId"), ` IN (SELECT user_id FROM group_member WHERE group_member.group_id = problem_set.group_id)))`)
	} else {
		sqlString += ` WHERE 1 = 1`
	}
	sqlString += fmt.Sprintf(`AND (query @@ document OR similarity > 0) ORDER BY rank_name, rank_description, similarity DESC NULLS LAST LIMIT $2 OFFSET $3`)
	var problemSets []model.ProblemSet
	if err := global.Database.Select(&problemSets, sqlString, request.Keyword, request.Limit, request.Offset); err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
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

// SearchGroups godoc
// @Schemes http
// @Description 搜索小组
// @Tags Search
// @Param search query SearchRequest true "搜索信息"
// @Success 200 {object} AllGroupResponse "小组列表"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /search/group [get]
// @Security ApiKeyAuth
func SearchGroups(c *gin.Context) {
	var request SearchRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.String(http.StatusBadRequest, "请求错误")
		return
	}
	sqlString := fmt.Sprintf(`SELECT
			"group".*
		FROM
			"group", 
			to_tsvector("group".name || "group".description) document,
			to_tsquery($1) query,
			NULLIF(ts_rank(to_tsvector("group".name), query), 0) rank_name,
			NULLIF(ts_rank(to_tsvector("group".description), query), 0) rank_description,
			SIMILARITY($1, "group".name || "group".description) similarity
		WHERE
		    query @@ document OR similarity > 0
		ORDER BY rank_name, rank_description, similarity DESC NULLS LAST LIMIT $2 OFFSET $3`)
	var groups []model.Group
	if err := global.Database.Select(&groups, sqlString, request.Keyword, request.Limit, request.Offset); err != nil {
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
		var area string
		sqlString = `SELECT name FROM area WHERE id = $1`
		if err := global.Database.Get(&area, sqlString, group.AreaId); err != nil {
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
			AreaName:    area,
			AvatarURL:   group.AvatarURL,
		})
	}
	c.JSON(http.StatusOK, AllGroupResponse{
		TotalCount: len(groupResponses),
		Group:      groupResponses,
	})
}

// SearchNotes godoc
// @Schemes http
// @Description 搜索笔记
// @Tags Search
// @Param search query SearchRequest true "搜索信息"
// @Success 200 {object} AllNoteResponse "笔记列表"
// @Failure 400 {string} string "请求解析失败"
// @Failure default {string} string "服务器错误"
// @Router /search/note [get]
// @Security ApiKeyAuth
func SearchNotes(c *gin.Context) {
	var request SearchRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.String(http.StatusBadRequest, "请求错误")
		return
	}
	sqlString := fmt.Sprintf(`SELECT
			note.*
		FROM
			note, 
			to_tsvector(note.title || note.content) document,
			to_tsquery($1) query,
			NULLIF(ts_rank(to_tsvector(note.title), query), 0) rank_name,
			NULLIF(ts_rank(to_tsvector(note.content), query), 0) rank_description,
			SIMILARITY($1, note.title || note.content) similarity`)
	role, _ := c.Get("Role")
	if role == global.GUEST {
		sqlString += ` WHERE is_public = true`
	} else if role == global.USER {
		sqlString += fmt.Sprint(` WHERE (is_public = true OR (user_id = `, c.GetInt("UserId"), `))`)
	} else {
		sqlString += ` WHERE 1 = 1`
	}
	sqlString += fmt.Sprintf(`AND (query @@ document OR similarity > 0) ORDER BY rank_name, rank_description, similarity DESC NULLS LAST LIMIT $2 OFFSET $3`)
	var notes []model.Note
	if err := global.Database.Select(&notes, sqlString, request.Keyword, request.Limit, request.Offset); err != nil {
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
