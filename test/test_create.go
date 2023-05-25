package test

import (
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func TestCreateGroup(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[2].Name,
		Password: initUser[2].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")

	// 再创建小组
	var result string
	var fill interface{}
	code = Post("/group/create", res.Token, &api.GroupCreateRequest{
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

func TestCreateNote(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[1].Name,
		Password: initUser[1].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")

	// 再创建笔记
	var result string
	var fill interface{}
	code = Post("/note/create", res.Token, &api.NoteCreateRequest{
		Title:    "test note",
		Content:  "test note",
		IsPublic: true,
	}, &result)
	assert.Equal(t, code, http.StatusOK)

	// 再查询
	var note api.AllNoteResponse
	code = Get("/note/all", res.Token, make(map[string][]string), &note)
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, len(note.Notes), 6)

	// 再更新
	newTitle := "new title"
	code = Put("/note/update", res.Token, &api.NoteUpdateRequest{
		ID:       initNote[1].ID,
		Title:    &newTitle,
		Content:  &initNote[1].Content,
		IsPublic: &initNote[1].IsPublic,
	}, &result)
	assert.Equal(t, code, http.StatusOK)

	// 再删除笔记
	code = Delete("/note/delete/9", res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)
}
