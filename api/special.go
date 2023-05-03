package api

import "github.com/gin-gonic/gin"

// This file will implement some special APIs, which is not in our expectation.
// I do not know why we need it, but front-end requires it.

// GetWrongProblemSet godoc
// @Schemes http
// @Description 获取错题集
// @Tags special
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"

// @Failure default {string} string "服务器错误"
// @Router /special/wrong_problem_set [get]
// @Security ApiKeyAuth
func GetWrongProblemSet(c *gin.Context) {

}

// GetFavoriteProblemSet godoc
// @Schemes http
// @Description 获取收藏题集
// @Tags special
// @Param offset query int false "偏移量"
// @Param limit query int false "数量"

// @Failure default {string} string "服务器错误"
// @Router /special/favorite_problem_set [get]
func GetFavoriteProblemSet(c *gin.Context) {

}
