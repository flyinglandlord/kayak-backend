package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Ping godoc
// @Schemes http
// @Description 测试服务器是否正常运行
// @Tags Availability
// @Success 200 {string} string "pong"
// @Failure default {string} string "服务器错误"
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
