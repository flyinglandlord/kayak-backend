package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testGetGroups(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var res interface{}
	code = Get(fmt.Sprintf("/group/all"), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/group/all"), loginRes.Token, map[string][]string{
		"id":       {"1"},
		"user_id":  {"1"},
		"owner_id": {"1"},
		"area_id":  {"1"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)
}

func testCreateGroup(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUsers[1]
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	_loginRes := api.LoginResponse{}
	_user := randomUsers[2]
	code = Post("/login", "", &api.LoginInfo{
		UserName: _user.Name,
		Password: _user.Password,
	}, &_loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, _loginRes.Token, "")

	var groupRes api.GroupResponse
	code = Post("/group/create", loginRes.Token, &api.GroupCreateRequest{
		Name:        "test",
		Description: "test",
	}, &groupRes)
	assert.Equal(t, code, http.StatusOK)

	var res interface{}
	code = Get(fmt.Sprintf("/group/invitation/%d", 100), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Get(fmt.Sprintf("/group/invitation/%d", groupRes.Id), _loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/group/invitation/%d", groupRes.Id), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/group/all_user/%d", 100), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Get(fmt.Sprintf("/group/all_user/%d", groupRes.Id), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Put(fmt.Sprintf("/group/update/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Put(fmt.Sprintf("/group/update/%d", groupRes.Id), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Put(fmt.Sprintf("/group/update/%d", groupRes.Id), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/group/apply"), loginRes.Token, &api.ApplyToJoinGroupRequest{
		GroupId: 100,
	}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/group/apply"), loginRes.Token, &api.ApplyToJoinGroupRequest{
		GroupId: groupRes.Id,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/group/apply"), _loginRes.Token, &api.ApplyToJoinGroupRequest{
		GroupId: groupRes.Id,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/group/apply"), _loginRes.Token, &api.ApplyToJoinGroupRequest{
		GroupId: groupRes.Id,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/group/application/%d?status=0", 100), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Get(fmt.Sprintf("/group/application/%d?status=0", groupRes.Id), _loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/group/application/%d?status=4", groupRes.Id), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	var applyRes api.GroupApplicationResponse
	code = Get(fmt.Sprintf("/group/application/%d?status=0&offset=0&limit=10", groupRes.Id), loginRes.Token, map[string][]string{}, &applyRes)
	assert.Equal(t, code, http.StatusOK)

	code = Put(fmt.Sprintf("/group/application"), loginRes.Token, &api.HandleGroupApplicationRequest{
		ApplicationId: 100,
		Status:        1,
	}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Put(fmt.Sprintf("/group/application"), _loginRes.Token, &api.HandleGroupApplicationRequest{
		ApplicationId: applyRes.Applications[0].Application.ID,
		Status:        1,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Put(fmt.Sprintf("/group/application"), loginRes.Token, &api.HandleGroupApplicationRequest{
		ApplicationId: applyRes.Applications[0].Application.ID,
		Status:        1,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Put(fmt.Sprintf("/group/application"), loginRes.Token, &api.HandleGroupApplicationRequest{
		ApplicationId: applyRes.Applications[0].Application.ID,
		Status:        1,
	}, &res)
	assert.Equal(t, code, http.StatusBadRequest)

	code = Delete(fmt.Sprintf("/group/quit/%d", 100), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/group/quit/%d", groupRes.Id), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/group/quit/%d", groupRes.Id), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/group/apply"), _loginRes.Token, &api.ApplyToJoinGroupRequest{
		GroupId: groupRes.Id,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/group/application/%d?status=0&offset=0&limit=10", groupRes.Id), loginRes.Token, map[string][]string{}, &applyRes)
	assert.Equal(t, code, http.StatusOK)

	code = Put(fmt.Sprintf("/group/application"), loginRes.Token, &api.HandleGroupApplicationRequest{
		ApplicationId: applyRes.Applications[0].Application.ID,
		Status:        1,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/group/remove/%d?user_id=%d", 100, 100), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/group/remove/%d?user_id=%d", groupRes.Id, 100), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/group/remove/%d?user_id=%d", groupRes.Id, 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/group/delete/%d", 100), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/group/delete/%d", groupRes.Id), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/group/delete/%d", groupRes.Id), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)
}
