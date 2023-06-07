package test

import (
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testLogin(t *testing.T) {
	var str interface{}
	var res api.LoginResponse
	user := randomUser()
	code := Post("/login", "", &struct{}{}, &str)
	assert.Equal(t, code, http.StatusBadRequest)
	code = Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password + user.Password,
	}, &str)
	assert.Equal(t, code, http.StatusBadRequest)
	code = Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
}

func testRegister(t *testing.T) {
	var res interface{}
	user := randomUser()
	code := Post("/register", "", &api.RegisterInfo{
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
		VCode:    "0",
	}, &res)
	assert.Equal(t, code, 409)
	code = Post("/register", "", &api.RegisterInfo{
		Name:     user.Name + user.Name,
		Password: user.Password,
		Email:    user.Email,
		VCode:    "0",
	}, &res)
	assert.Equal(t, code, http.StatusBadRequest)
}

func testChangePassword(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
	// 再重置密码
	var result interface{}
	code = Post("/change-password", res.Token, &api.RegisterResponse{
		OldPassword: user.Password + user.Password,
		NewPassword: user.Password,
	}, &result)
	assert.Equal(t, code, http.StatusBadRequest)
	code = Post("/change-password", res.Token, &api.RegisterResponse{
		OldPassword: user.Password,
		NewPassword: user.Password + user.Password,
	}, &result)
	assert.Equal(t, code, http.StatusOK)
	code = Post("/change-password", res.Token, &api.RegisterResponse{
		OldPassword: user.Password + user.Password,
		NewPassword: user.Password,
	}, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testResetPassword(t *testing.T) {
	var result interface{}
	user := randomUser()
	code := Post("/reset-password", "", &api.ResetPasswordInfo{
		UserName:    user.Name + user.Name,
		VerifyCode:  "0",
		NewPassword: user.Password,
	}, &result)
	assert.Equal(t, code, 409)
	code = Post("/reset-password", "", &api.ResetPasswordInfo{
		UserName:    user.Name,
		VerifyCode:  "0",
		NewPassword: user.Password,
	}, &result)
	assert.Equal(t, code, http.StatusBadRequest)
}

func testLogout(t *testing.T) {
	res := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
	var result string
	code = Get("/logout", res.Token, nil, &result)
	assert.Equal(t, code, http.StatusOK)
}
