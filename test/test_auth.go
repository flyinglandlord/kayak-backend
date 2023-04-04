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
