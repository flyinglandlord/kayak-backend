package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testAddDiscussionReview(t *testing.T) {
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

	var discussionRes api.DiscussionResponse
	code = Post("/discussion/create", loginRes.Token, &api.DiscussionCreateRequest{
		Title:    "test",
		Content:  "test",
		GroupId:  groupRes.Id,
		IsPublic: false,
	}, &discussionRes)
	assert.Equal(t, code, http.StatusOK)

	var discussionReviewRes api.DiscussionReviewResponse
	code = Post("/discussion_review/add", loginRes.Token, &api.DiscussionReviewCreateRequest{
		Title:        "test",
		Content:      "test",
		DiscussionId: 100,
	}, &discussionReviewRes)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post("/discussion_review/add", _loginRes.Token, &api.DiscussionReviewCreateRequest{
		Title:        "test",
		Content:      "test",
		DiscussionId: discussionRes.ID,
	}, &discussionReviewRes)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post("/discussion_review/add", loginRes.Token, &api.DiscussionReviewCreateRequest{
		Title:        "test",
		Content:      "test",
		DiscussionId: discussionRes.ID,
	}, &discussionReviewRes)
	assert.Equal(t, code, http.StatusOK)

	var res interface{}
	code = Get(fmt.Sprintf("/discussion_review/get?discussion_id=100"), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Get(fmt.Sprintf("/discussion_review/get?discussion_id=%d", discussionRes.ID), _loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/discussion_review/get?discussion_id=%d", discussionRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion_review/like/100"), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/discussion_review/like/%d", discussionReviewRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/discussion_review/like/%d", discussionReviewRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion_review/like/%d", discussionReviewRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/discussion_review/unlike/100"), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/discussion_review/unlike/%d", discussionReviewRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/discussion_review/remove/100"), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/discussion_review/remove/%d", discussionReviewRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/discussion_review/remove/%d", discussionReviewRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)
}
