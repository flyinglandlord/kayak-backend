package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
	"time"
)

func testUserInfo(t *testing.T) {
	loginRes := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[0].Name,
		Password: initUser[0].Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	token := loginRes.Token
	userInfoRes := api.UserInfoResponse{}
	code = Get("/user/info", token, map[string][]string{}, &userInfoRes)
	fmt.Println(time.Now().Local())
	fmt.Println(userInfoRes.CreateAt)
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, userInfoRes.UserId, initUser[0].ID)
	assert.Equal(t, userInfoRes.UserName, initUser[0].Name)
	assert.Equal(t, userInfoRes.Email, initUser[0].Email)
	assert.Equal(t, userInfoRes.Phone, initUser[0].Phone)
}
