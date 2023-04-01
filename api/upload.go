package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/viper"
	"kayak-backend/global"
	"net/http"
	"path"
	"strings"
)

func SanitizeFilename(s string) string {
	r := path.Clean(path.Base(s))
	blacklist := []rune("" +
		// S3 Avoids
		"\\{}^%`[]'\"<>~#|" +
		// S3 Requires Handling
		"&$@=;/:+ ,?" +
		"\x7f\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f" +
		"\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f" +
		// Conflicts with Markdown
		"()")
	for _, ch := range blacklist {
		r = strings.ReplaceAll(r, string(ch), "_")
	}
	return r
}

func DoUploadPublic(c *gin.Context) (int, string) {
	UserId := c.GetInt("UserId")
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		return http.StatusBadRequest, ""
	}

	// construct objectPath
	filename := SanitizeFilename(fileHeader.Filename)
	randomId := uuid.New().String()
	objectPath := fmt.Sprintf("%d/%s/%s", UserId, randomId, filename)

	_, err = global.MinioClient.PutObject(
		"public",
		objectPath,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return http.StatusBadGateway, "上传失败"
	}

	url := viper.GetString("S3PublicBucketRoute") + "/" + objectPath
	return http.StatusOK, url
}

// UploadPublicFile godoc
// @Schemes http
// @Description 上传公开文件
// @Param file formData file true "文件"
// @Success 200 {string} string "文件 URL"
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "请求被禁止"
// @Failure default {string} string "服务器错误"
// @Router /upload/public [post]
// @Security ApiKeyAuth
func UploadPublicFile(c *gin.Context) {
	status, url := DoUploadPublic(c)
	c.JSON(status, gin.H{
		"url": "/public" + url,
	})
}

// UploadAvatar godoc
// @Schemes http
// @Description 上传用户头像
// @Param file formData file true "头像"
// @Success 200 {string} string
// @Failure 400 {string} string "请求解析失败"
// @Failure 403 {string} string "请求被禁止"
// @Failure default {string} string "服务器错误"
// @Router /upload/avatar [post]
// @Security ApiKeyAuth
func UploadAvatar(c *gin.Context) {
	status, url := DoUploadPublic(c)
	UserId := c.GetInt("UserId")
	if status == http.StatusOK {
		// 更新数据库
		sqlString := `UPDATE "user" SET avatar_url = $1 WHERE id = $2`
		if _, err := global.Database.Exec(sqlString, url, UserId); err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		c.String(http.StatusOK, "头像设置成功")
		return
	}
	c.String(status, "头像设置失败")
}
