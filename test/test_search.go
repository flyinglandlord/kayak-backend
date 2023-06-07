package test

import (
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testSearchProblemSets(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var result api.AllProblemSetResponse
	code = Post("/search/problem_set", loginRes.Token, &api.SearchRequest{
		Keyword: "problemSet",
		Limit:   1,
		Offset:  0,
	}, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testSearchGroups(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var result api.AllGroupResponse
	code = Post("/search/group", loginRes.Token, &api.SearchRequest{
		Keyword: "group",
		Limit:   1,
		Offset:  0,
	}, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testSearchNotes(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var result api.AllNoteResponse
	code = Post("/search/note", loginRes.Token, &api.SearchRequest{
		Keyword: "note",
		Limit:   1,
		Offset:  0,
	}, &result)
	assert.Equal(t, code, http.StatusOK)
}
