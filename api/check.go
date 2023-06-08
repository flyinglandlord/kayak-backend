package api

import (
	"github.com/gin-gonic/gin"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
)

// CheckProblemSetWriteAuth godoc
// @Schemes http
// @Description 检查用户是否有权限修改题集
// @Tags Check
// @Param id path int true "题集ID"
// @Success 200 {string} string "用户有权限修改题集"
// @Failure 403 {string} string "用户没有权限修改题集"
// @Failure 404 {string} string "题集不存在"
// @Failure default {string} string "服务器错误"
// @Router /check/problem_set/{id} [get]
// @Security ApiKeyAuth
func CheckProblemSetWriteAuth(c *gin.Context) {
	var problemSet model.ProblemSet
	sqlString := `SELECT * FROM problem_set WHERE id = $1`
	if err := global.Database.Get(&problemSet, sqlString, c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "题集不存在")
		return
	}
	role, _ := c.Get("Role")
	if problemSet.GroupId == 0 {
		if role != global.ADMIN && problemSet.UserId != c.GetInt("UserId") {
			c.String(http.StatusForbidden, "用户没有权限修改题集")
			return
		}
	} else {
		sqlString = `SELECT count(*) FROM group_member WHERE group_id = $1 AND user_id = $2`
		var count int
		if err := global.Database.Get(&count, sqlString, problemSet.GroupId, c.GetInt("UserId")); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		if role != global.ADMIN && count == 0 {
			c.String(http.StatusForbidden, "用户没有权限修改题集")
			return
		}
	}
	c.String(http.StatusOK, "用户有权限修改题集")
}
