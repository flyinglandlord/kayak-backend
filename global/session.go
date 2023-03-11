package global

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type Role int

const TokenHeader = "X-Token"

const (
	GUEST Role = iota
	USER
	ADMIN
)

type Session struct {
	Role   Role `json:"role"`
	UserId int  `json:"user_id"`
}

func GetSessionByToken(c context.Context, token string) []byte {
	r, err := Redis.Get(c, token).Bytes()
	if err != nil {
		return nil
	}
	return r
}

func CreateSession(c context.Context, session *Session) (string, error) {
	token := fmt.Sprintf("%s@%d", uuid.New().String(), session.UserId)
	bytes, err := json.Marshal(*session)
	if err != nil {
		return "", err
	}
	err = Redis.Set(c, token, bytes, 0).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func DeleteSession(c context.Context, token string) error {
	err := Redis.Del(c, token).Err()
	return err
}

func Authenticate(c *gin.Context) {
	token := c.Request.Header.Get(TokenHeader)
	if token == "" {
		c.Set("Role", GUEST)
		c.Next()
		return
	}
	sessionInfo := GetSessionByToken(c, token)
	if sessionInfo == nil {
		c.Set("Role", GUEST)
		c.Next()
		return
	}
	var session Session
	fmt.Println(string(sessionInfo))
	err := json.Unmarshal(sessionInfo, &session)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器错误")
		c.Abort()
		return
	}
	c.Set("Role", session.Role)
	c.Set("UserId", session.UserId)
	c.Next()
}

func CheckAuth(c *gin.Context) {
	role, ok := c.Get("Role")
	if !ok || role == GUEST {
		c.String(http.StatusUnauthorized, "未登录")
		c.Abort()
		return
	}
	c.Next()
}
