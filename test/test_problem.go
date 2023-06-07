package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testGetChoiceProblems(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var res interface{}
	code = Get(fmt.Sprintf("/problem/choice/all"), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/choice/all"), loginRes.Token, map[string][]string{
		"id":          {"1"},
		"user_id":     {"1"},
		"is_favorite": {"true"},
		"is_wrong":    {"true"},
		"offset":      {"0"},
		"limit":       {"1"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/choice/all"), loginRes.Token, map[string][]string{
		"is_wrong": {"true"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/choice/all"), loginRes.Token, map[string][]string{
		"is_wrong": {"false"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)
}

func testGetBlankProblems(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var res interface{}
	code = Get(fmt.Sprintf("/problem/blank/all"), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/blank/all"), loginRes.Token, map[string][]string{
		"id":          {"1"},
		"user_id":     {"1"},
		"is_favorite": {"true"},
		"is_wrong":    {"true"},
		"offset":      {"0"},
		"limit":       {"1"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/blank/all"), loginRes.Token, map[string][]string{
		"is_wrong": {"true"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/blank/all"), loginRes.Token, map[string][]string{
		"is_wrong": {"false"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)
}

func testGetJudgeProblems(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var res interface{}
	code = Get(fmt.Sprintf("/problem/judge/all"), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/judge/all"), loginRes.Token, map[string][]string{
		"id":          {"1"},
		"user_id":     {"1"},
		"is_favorite": {"true"},
		"is_wrong":    {"true"},
		"offset":      {"0"},
		"limit":       {"1"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/judge/all"), loginRes.Token, map[string][]string{
		"is_wrong": {"true"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/judge/all"), loginRes.Token, map[string][]string{
		"is_wrong": {"false"},
	}, &res)
	assert.Equal(t, code, http.StatusOK)
}
