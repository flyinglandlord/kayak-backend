package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testGetWrongProblemSet(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var problemSetRes api.ProblemSetResponse
	code = Post("/problem_set/create", loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    false,
	}, &problemSetRes)
	assert.Equal(t, code, http.StatusOK)
	var problemRes api.JudgeProblemResponse
	code = Post("/problem/judge/create", loginRes.Token, &api.JudgeProblemCreateRequest{
		Description: "test",
		IsPublic:    true,
		IsCorrect:   true,
	}, &problemRes)
	assert.Equal(t, code, http.StatusOK)
	var res interface{}
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", problemSetRes.ID, problemRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)
	code = Post(fmt.Sprintf("/wrong_record/create/%d", problemRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	var result api.AllWrongProblemSet
	code = Get("/special/wrong_problem_set", loginRes.Token, map[string][]string{
		"offset": {"0"},
		"limit":  {"50"},
	}, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testGetFavoriteProblemSet(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var problemSetRes api.ProblemSetResponse
	code = Post("/problem_set/create", loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    false,
	}, &problemSetRes)
	assert.Equal(t, code, http.StatusOK)
	var problemRes api.JudgeProblemResponse
	code = Post("/problem/judge/create", loginRes.Token, &api.JudgeProblemCreateRequest{
		Description: "test",
		IsPublic:    true,
		IsCorrect:   true,
	}, &problemRes)
	assert.Equal(t, code, http.StatusOK)
	var res interface{}
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", problemSetRes.ID, problemRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)
	code = Post(fmt.Sprintf("/problem/favorite/%d", problemRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	var result api.AllFavoriteProblemSet
	code = Get("/special/favorite_problem_set", loginRes.Token, map[string][]string{
		"offset": {"0"},
		"limit":  {"50"},
	}, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testGetFeaturedProblemSet(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var result api.AllProblemSetResponse
	code = Get("/special/featured_problem_set", loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testGetFeaturedNote(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var result api.AllNoteResponse
	code = Get("/special/featured_note", loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testGetFeaturedGroup(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var result api.AllGroupResponse
	code = Get("/special/featured_group", loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusOK)
}
