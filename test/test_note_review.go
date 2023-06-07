package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testAddNoteReview(t *testing.T) {
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

	var noteRes api.NoteResponse
	code = Post("/note/create", loginRes.Token, &api.NoteCreateRequest{
		Title:    "test",
		Content:  "test",
		IsPublic: false,
		Problems: []int{},
	}, &noteRes)
	assert.Equal(t, code, http.StatusOK)

	var res interface{}
	code = Post("/note_review/add", loginRes.Token, &api.NoteReviewCreateRequest{
		Title:   "test",
		Content: "test",
		NoteId:  100,
	}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post("/note_review/add", _loginRes.Token, &api.NoteReviewCreateRequest{
		Title:   "test",
		Content: "test",
		NoteId:  noteRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	var noteReviewRes api.NoteReviewResponse
	code = Post("/note_review/add", loginRes.Token, &api.NoteReviewCreateRequest{
		Title:   "test",
		Content: "test",
		NoteId:  noteRes.ID,
	}, &noteReviewRes)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/note_review/get?note_id=%d", 100), loginRes.Token, map[string][]string{}, &noteReviewRes)
	assert.Equal(t, code, http.StatusNotFound)

	code = Get(fmt.Sprintf("/note_review/get?note_id=%d", noteRes.ID), _loginRes.Token, map[string][]string{}, &noteReviewRes)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/note_review/get?note_id=%d", noteRes.ID), loginRes.Token, map[string][]string{}, &noteReviewRes)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/note_review/like/%d", 100), loginRes.Token, map[string][]string{}, &noteReviewRes)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/note_review/like/%d", noteReviewRes.ID), loginRes.Token, map[string][]string{}, &noteReviewRes)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/note_review/unlike/%d", 100), loginRes.Token, map[string][]string{}, &noteReviewRes)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/note_review/unlike/%d", noteReviewRes.ID), loginRes.Token, map[string][]string{}, &noteReviewRes)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/note_review/remove/%d", 100), loginRes.Token, struct{}{}, &noteReviewRes)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/note_review/remove/%d", noteReviewRes.ID), _loginRes.Token, struct{}{}, &noteReviewRes)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/note_review/remove/%d", noteReviewRes.ID), loginRes.Token, struct{}{}, &noteReviewRes)
	assert.Equal(t, code, http.StatusOK)
}
