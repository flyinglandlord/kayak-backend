package test

import (
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"strconv"
	"testing"
)

func testFavoriteNote(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[2].Name,
		Password: initUser[2].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
	// 再收藏笔记
	var result string
	var fill interface{}
	code = Post("/note/favorite/"+strconv.Itoa(initNote[6].ID), res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)

	// 再取消收藏笔记
	code = Delete("/note/unfavorite/"+strconv.Itoa(initNote[6].ID), res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testLikeNoteReview(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[2].Name,
		Password: initUser[2].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
	// 再收藏笔记评论
	var result string
	var fill interface{}
	code = Post("/note_review/like/"+strconv.Itoa(initNoteReview[11].ID), res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)

	// 再取消收藏笔记评论
	code = Post("/note_review/unlike/"+strconv.Itoa(initNoteReview[11].ID), res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testLikeNote(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[2].Name,
		Password: initUser[2].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
	// 再收藏笔记
	var result string
	var fill interface{}
	code = Post("/note/like/"+strconv.Itoa(initNote[6].ID), res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)

	// 再取消收藏笔记
	code = Post("/note/unlike/"+strconv.Itoa(initNote[6].ID), res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testFavoriteProblem(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[2].Name,
		Password: initUser[2].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
	// 再收藏题目
	var result string
	var fill interface{}
	code = Post("/problem/favorite/"+strconv.Itoa(initProblemType[1].ID), res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)

	// 再取消收藏题目
	code = Delete("/problem/unfavorite/"+strconv.Itoa(initProblemType[1].ID), res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testFavoriteProblemSet(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[2].Name,
		Password: initUser[2].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")
	// 再收藏题目集
	var result string
	var fill interface{}
	code = Post("/problem_set/favorite/"+strconv.Itoa(initProblemSet[6].ID), res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)

	// 再取消收藏题目集
	code = Delete("/problem_set/unfavorite/"+strconv.Itoa(initProblemSet[6].ID), res.Token, &fill, &result)
	assert.Equal(t, code, http.StatusOK)
}
