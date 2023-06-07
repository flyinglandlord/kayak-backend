package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testCreateWrongRecord(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUsers[0]
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	_loginRes := api.LoginResponse{}
	_user := randomUsers[1]
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

	var res interface{}
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", problemSetRes.ID, blankRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/wrong_record/create/100"), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/wrong_record/create/%d", blankRes.ID), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/wrong_record/create/%d", blankRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/wrong_record/get/%d", blankRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/wrong_record/delete/%d", blankRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)
}
