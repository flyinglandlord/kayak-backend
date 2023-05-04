package test

import (
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"strconv"
	"testing"
)

func TestProblemAnswer(t *testing.T) {
	// 先登录
	res := api.LoginResponse{}
	code := Post("/login", "", &api.LoginInfo{
		UserName: initUser[2].Name,
		Password: initUser[2].Password,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, res.Token, "")

	// 获取选择题
	var list []api.ChoiceProblemResponse
	code = Get("/problem/choice/all", res.Token, make(map[string][]string), &list)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, list, 2)

	// 获取选择题答案
	var answer api.ChoiceProblemAnswerResponse
	code = Get("/problem/choice/answer/"+strconv.Itoa(initProblemType[0].ID), res.Token, make(map[string][]string), &answer)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, answer.ChoiceProblemAnswer[0].IsCorrect, false)

	// 获取判断题
	var list2 []api.BlankProblemResponse
	code = Get("/problem/blank/all", res.Token, make(map[string][]string), &list2)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, list2, 2)

	// 获取判断题答案
	var answer2 api.BlankProblemAnswerResponse
	code = Get("/problem/blank/answer/"+strconv.Itoa(initProblemType[1].ID), res.Token, make(map[string][]string), &answer2)
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, answer2.Answer, "problem2_answer")
}
