package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testGetNotes(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var res interface{}
	code = Get(fmt.Sprintf("/note/all"), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/note/all"), loginRes.Token, map[string][]string{
		"id":           {"1"},
		"user_id":      {"1"},
		"is_liked":     {"true"},
		"is_favorite":  {"true"},
		"offset":       {"0"},
		"limit":        {"1"},
		"sort_by_like": {"true"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/note/all"), loginRes.Token, map[string][]string{
		"is_liked":     {"false"},
		"is_favorite":  {"false"},
		"sort_by_like": {"false"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)
}

func testCreateNote(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUsers[4]
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	_loginRes := api.LoginResponse{}
	_user := randomUsers[0]
	code = Post("/login", "", &api.LoginInfo{
		UserName: _user.Name,
		Password: _user.Password,
	}, &_loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, _loginRes.Token, "")

	var problemSetRes api.ProblemSetResponse
	code = Post("/problem_set/create", loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    true,
	}, &problemSetRes)
	assert.Equal(t, code, http.StatusOK)

	var blankRes api.BlankProblemResponse
	code = Post("/problem/blank/create", loginRes.Token, &api.BlankProblemCreateRequest{
		Description:   "test",
		IsPublic:      false,
		Answer:        "test",
		AnswerExplain: "test",
	}, &blankRes)
	assert.Equal(t, code, http.StatusOK)

	var _blankRes api.BlankProblemResponse
	code = Post("/problem/blank/create", loginRes.Token, &api.BlankProblemCreateRequest{
		Description:   "test",
		IsPublic:      false,
		Answer:        "test",
		AnswerExplain: "test",
	}, &_blankRes)
	assert.Equal(t, code, http.StatusOK)

	var res interface{}
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", problemSetRes.ID, blankRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", problemSetRes.ID, _blankRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	var noteRes api.NoteResponse
	code = Post("/note/create", loginRes.Token, &api.NoteCreateRequest{
		Title:    "test",
		Content:  "test",
		IsPublic: false,
		Problems: []int{blankRes.ID},
	}, &noteRes)
	assert.Equal(t, code, http.StatusOK)

	code = Put("/note/update", loginRes.Token, &api.NoteUpdateRequest{
		ID: 100,
	}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Put("/note/update", _loginRes.Token, &api.NoteUpdateRequest{
		ID: noteRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Put("/note/update", loginRes.Token, &api.NoteUpdateRequest{
		ID: noteRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/note/like/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/note/like/%d", noteRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/note/like/%d", noteRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/note/unlike/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/note/unlike/%d", noteRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/note/favorite/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/note/favorite/%d", noteRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/note/favorite/%d", noteRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/note/unfavorite/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/note/unfavorite/%d", noteRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/note/add_problem/%d?problem_id=%d", 100, 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/note/add_problem/%d?problem_id=%d", noteRes.ID, 100), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/note/add_problem/%d?problem_id=%d", noteRes.ID, 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/note/add_problem/%d?problem_id=%d", noteRes.ID, _blankRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/note/remove_problem/%d?problem_id=%d", 100, 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/note/remove_problem/%d?problem_id=%d", noteRes.ID, 100), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/note/remove_problem/%d?problem_id=%d", noteRes.ID, 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/note/remove_problem/%d?problem_id=%d", noteRes.ID, _blankRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/note/problem_list/%d?offset=0&limit=1", noteRes.ID), _loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/note/problem_list/%d?offset=0&limit=1", noteRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/note/problem_list/%d?offset=0&limit=1", 100), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/note/remove_problem/%d?problem_id=%d", noteRes.ID, blankRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/note/problem_list/%d?offset=0&limit=1", noteRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/note/delete/%d", 100), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/note/delete/%d", noteRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/note/delete/%d", noteRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)
}
