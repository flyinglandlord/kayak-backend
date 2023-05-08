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
// @Tags special
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
// @Tags special
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
