package test

import (
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testCreateDiscussion(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[0].Name,
		Password: initUser[0].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")

	// 再创建讨论
	var result string
	var fill interface{}
	code = Post("/discussion/create", res.Token, &api.GroupCreateRequest{
		Name:        "test group",
		Description: "test group",
		AreaId:      nil,
	}, &result)
	assert.Equal(t, code, http.StatusOK)

	// 再查询
	var group api.AllGroupResponse
	code = Get("/group/all", res.Token, make(map[string][]string), &group)
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, group.Group[0].Name, "test group")
	assert.Equal(t, group.Group[0].Description, "test group")
	assert.Equal(t, len(group.Group), 1)

	// 再删除小组
	code = Delete("/group/delete/1", res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)
}
