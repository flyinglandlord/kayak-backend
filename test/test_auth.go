package test

import (
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testLogin(t *testing.T) {
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[0].Name,
		Password: initUser[0].Password + initUser[0].Password,
	}, &res)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, res.Token, "")
	code = Post("/login", "", &api.LoginInfo{
		UserName: initUser[0].Name,
		Password: initUser[0].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
}

func testRegister(t *testing.T) {
	var res string
	code := Post("/register", "", &api.RegisterInfo{
		Name:     initUser[1].Name,
		Password: initUser[1].Password,
	}, &res)
	assert.Equal(t, code, 409)
	//assert.Equal(t, res, "用户名已存在")
	code = Post("/register", "", &api.RegisterInfo{
		Name:     "initUser[1].Name + initUser[1].Name",
		Password: "initUser[1].Password,",
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	//assert.NotEqual(t, res, "注册成功")
}

func testChangePassword(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[2].Name,
		Password: initUser[2].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
	// 再重置密码
	var result string
	code = Post("/change-password", res.Token, &api.RegisterResponse{
		OldPassword: initUser[2].Password + initUser[2].Password,
		NewPassword: initUser[2].Password,
	}, &result)
	assert.Equal(t, code, http.StatusBadRequest)
	//assert.Equal(t, result, "旧密码错误")
	code = Post("/change-password", res.Token, &api.RegisterResponse{
		OldPassword: initUser[2].Password,
		NewPassword: initUser[2].Password + initUser[2].Password,
	}, &result)
	assert.Equal(t, code, http.StatusOK)
	//assert.Equal(t, result, "修改成功")
	code = Post("/change-password", res.Token, &api.RegisterResponse{
		OldPassword: initUser[2].Password + initUser[2].Password,
		NewPassword: initUser[2].Password,
	}, &result)
	assert.Equal(t, code, http.StatusOK)
	//assert.Equal(t, result, "修改成功")
}

func testLogout(t *testing.T) {
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[3].Name,
		Password: initUser[3].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
	var result string
	code = Get("/logout", res.Token, nil, &result)
	assert.Equal(t, code, http.StatusOK)
	//assert.Equal(t, result, "退出成功")
}
