package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"strconv"
	"testing"
)

func testCreateDiscussion(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUsers[2]
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	_loginRes := api.LoginResponse{}
	_user := randomUsers[3]
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
	code = Post("/discussion/create", loginRes.Token, &api.DiscussionCreateRequest{
		Title:    "test",
		Content:  "test",
		GroupId:  100,
		IsPublic: true,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	var discussionRes api.DiscussionResponse
	code = Post("/discussion/create", loginRes.Token, &api.DiscussionCreateRequest{
		Title:    "test",
		Content:  "test",
		GroupId:  groupRes.Id,
		IsPublic: false,
	}, &discussionRes)
	assert.Equal(t, code, http.StatusOK)

	var _discussionRes api.DiscussionResponse
	code = Post("/discussion/create", loginRes.Token, &api.DiscussionCreateRequest{
		Title:    "test",
		Content:  "test",
		GroupId:  groupRes.Id,
		IsPublic: true,
	}, &_discussionRes)
	assert.Equal(t, code, http.StatusOK)

	code = Put("/discussion/update", loginRes.Token, &api.DiscussionUpdateRequest{
		ID: 100,
	}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Put("/discussion/update", _loginRes.Token, &api.DiscussionUpdateRequest{
		ID: discussionRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Put("/discussion/update", loginRes.Token, &api.DiscussionUpdateRequest{
		ID: discussionRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion/like/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/discussion/like/%d", discussionRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/discussion/like/%d", _discussionRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/discussion/like/%d", discussionRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion/like/%d", discussionRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion/unlike/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/discussion/unlike/%d", discussionRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion/unlike/%d", discussionRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion/favorite/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/discussion/favorite/%d", discussionRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/discussion/favorite/%d", _discussionRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/discussion/favorite/%d", discussionRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion/favorite/%d", discussionRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion/unfavorite/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/discussion/unfavorite/%d", discussionRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion/unfavorite/%d", discussionRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/discussion/all"), loginRes.Token, map[string][]string{
		"group_id": {"100"},
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/discussion/all"), loginRes.Token, map[string][]string{
		"group_id": {strconv.Itoa(groupRes.Id)},
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/discussion/all"), loginRes.Token, map[string][]string{
		"group_id":     {strconv.Itoa(groupRes.Id)},
		"id":           {"1"},
		"user_id":      {"1"},
		"is_liked":     {"true"},
		"offset":       {"0"},
		"limit":        {"1"},
		"sort_by_like": {"true"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/discussion/all"), loginRes.Token, map[string][]string{
		"group_id":     {strconv.Itoa(groupRes.Id)},
		"is_liked":     {"false"},
		"offset":       {"0"},
		"limit":        {"1"},
		"sort_by_like": {"false"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/discussion/delete/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/discussion/delete/%d", discussionRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/discussion/delete/%d", discussionRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)
}
