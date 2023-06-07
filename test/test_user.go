package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testGetUserInfoById(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	token := loginRes.Token
	userInfoRes := api.UserInfoResponse{}
	_user := randomUser()
	code = Get(fmt.Sprintf("/user/info/%d", _user.ID), token, map[string][]string{}, &userInfoRes)
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, userInfoRes.UserId, _user.ID)
	assert.Equal(t, userInfoRes.NickName, _user.NickName)
}

func testGetUserInfo(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	token := loginRes.Token
	userInfoRes := api.UserInfoResponse{}
	code = Get("/user/info", token, map[string][]string{}, &userInfoRes)
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, userInfoRes.UserId, user.ID)
	assert.Equal(t, userInfoRes.UserName, user.Name)
	assert.Equal(t, userInfoRes.Email, user.Email)
}

func testUpdateUserInfo(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	token := loginRes.Token
	var res interface{}
	code = Put("/user/update", token, &api.UserInfoRequest{}, &res)
	assert.Equal(t, code, http.StatusOK)
}

func testGetUserWrongRecords(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	token := loginRes.Token
	var res api.AllWrongRecordResponse
	code = Get("/user/wrong_record", token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, res.TotalCount, 0)
}
