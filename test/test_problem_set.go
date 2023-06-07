package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"kayak-backend/api"
	"net/http"
	"testing"
)

func testGetProblemSets(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var result api.AllProblemSetResponse
	code = Get("/problem_set/all", loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusOK)
	code = Get("/problem_set/all", loginRes.Token, map[string][]string{
		"id":          {"1"},
		"user_id":     {"1"},
		"group_id":    {"1"},
		"is_favorite": {"true"},
		"contain":     {"1"},
		"area_id":     {"1"},
	}, &result)
	assert.Equal(t, code, http.StatusOK)
	code = Get("/problem_set/all", loginRes.Token, map[string][]string{
		"id":          {"1"},
		"user_id":     {"1"},
		"group_id":    {"0"},
		"is_favorite": {"false"},
		"contain":     {"1"},
		"area_id":     {"1"},
	}, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testCreateProblemSet(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var result interface{}
	code = Post("/problem_set/create", loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    true,
	}, &result)
	assert.Equal(t, code, http.StatusOK)
	groupId := -1
	code = Post("/problem_set/create", loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    true,
		GroupId:     &groupId,
	}, &result)
	assert.Equal(t, code, http.StatusForbidden)
}

func testUpdateProblemSet(t *testing.T) {
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

	var result api.ProblemSetResponse
	code = Post("/problem_set/create", loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    false,
	}, &result)
	assert.Equal(t, code, http.StatusOK)
	var problemRes api.JudgeProblemResponse
	code = Post("/problem/judge/create", loginRes.Token, &api.JudgeProblemCreateRequest{
		Description: "test",
		IsPublic:    true,
		IsCorrect:   true,
	}, &problemRes)
	assert.Equal(t, code, http.StatusOK)
	var res interface{}
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", result.ID, problemRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", result.ID, problemRes.ID), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)
	code = Get(fmt.Sprintf("/problem_set/all_problem/%d", result.ID), loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusOK)
	code = Get(fmt.Sprintf("/problem_set/all_problem/%d", result.ID), _loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusForbidden)
	code = Put("/problem_set/update", loginRes.Token, &api.ProblemSetUpdateRequest{
		ID: 100,
	}, &res)
	assert.Equal(t, code, http.StatusNotFound)
	code = Put("/problem_set/update", _loginRes.Token, &api.ProblemSetUpdateRequest{
		ID: result.ID,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)
	code = Put("/problem_set/update", loginRes.Token, &api.ProblemSetUpdateRequest{
		ID: result.ID,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	groupId := 100
	code = Put("/problem_set/update", loginRes.Token, &api.ProblemSetUpdateRequest{
		ID:      result.ID,
		GroupId: &groupId,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	var groupRes api.GroupResponse
	code = Post("/group/create", loginRes.Token, &api.GroupCreateRequest{
		Name:        "test",
		Description: "test",
	}, &groupRes)
	assert.Equal(t, code, http.StatusOK)
	code = Put("/problem_set/update", loginRes.Token, &api.ProblemSetUpdateRequest{
		ID:      result.ID,
		GroupId: &groupRes.Id,
	}, &res)
	assert.Equal(t, code, http.StatusOK)
	code = Put("/problem_set/update", _loginRes.Token, &api.ProblemSetUpdateRequest{
		ID:      result.ID,
		GroupId: &groupRes.Id,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)
	code = Get(fmt.Sprintf("/problem_set/all_problem/%d", result.ID), loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusOK)
	code = Get(fmt.Sprintf("/problem_set/all_problem/%d", result.ID), _loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusForbidden)
	var _problemRes api.JudgeProblemResponse
	code = Post("/problem/judge/create", _loginRes.Token, &api.JudgeProblemCreateRequest{
		Description: "test",
		IsPublic:    true,
		IsCorrect:   true,
	}, &_problemRes)
	assert.Equal(t, code, http.StatusOK)
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", result.ID, _problemRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", result.ID, _problemRes.ID), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)
}

func testGetProblemsInProblemSet(t *testing.T) {
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

	var result api.AllProblemResponse
	var problemSet api.ProblemSetResponse
	code = Post("/problem_set/create", loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    false,
	}, &problemSet)
	assert.Equal(t, code, http.StatusOK)
	code = Get(fmt.Sprintf("/problem_set/all_problem/100"), loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusNotFound)
	code = Get(fmt.Sprintf("/problem_set/all_problem/%d", problemSet.ID), loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusOK)
	code = Get(fmt.Sprintf("/problem_set/all_problem/%d", problemSet.ID), _loginRes.Token, map[string][]string{}, &result)
	assert.Equal(t, code, http.StatusForbidden)
	code = Get(fmt.Sprintf("/problem_set/all_problem/%d", problemSet.ID), loginRes.Token, map[string][]string{
		"is_favorite":     {"true"},
		"problem_type_id": {"0"},
		"is_wrong":        {"false"},
		"offset":          {"0"},
		"limit":           {"1"},
	}, &result)
	assert.Equal(t, code, http.StatusOK)
	code = Get(fmt.Sprintf("/problem_set/all_problem/%d", problemSet.ID), loginRes.Token, map[string][]string{
		"is_favorite":     {"false"},
		"problem_type_id": {"1"},
		"is_wrong":        {"true"},
		"offset":          {"0"},
		"limit":           {"1"},
	}, &result)
	assert.Equal(t, code, http.StatusOK)
}

func testAddProblemToProblemSet(t *testing.T) {
	loginRes := api.LoginResponse{}
	user := randomUser()
	code := Post("/login", "", &api.LoginInfo{
		UserName: user.Name,
		Password: user.Password,
	}, &loginRes)
	assert.Equal(t, code, http.StatusOK)
	assert.NotEqual(t, loginRes.Token, "")

	var result api.ProblemSetResponse
	code = Post("/problem_set/create", loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    false,
	}, &result)
	assert.Equal(t, code, http.StatusOK)
	var res interface{}
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", 100, 100), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", result.ID, 100), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)
}

func testMigrateProblemToProblemSet(t *testing.T) {
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

	var _groupRes api.GroupResponse
	code = Post("/group/create", _loginRes.Token, &api.GroupCreateRequest{
		Name:        "test",
		Description: "test",
	}, &_groupRes)
	assert.Equal(t, code, http.StatusOK)

	var problemSetRes api.ProblemSetResponse
	code = Post("/problem_set/create", loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    true,
	}, &problemSetRes)
	assert.Equal(t, code, http.StatusOK)

	var _problemSetRes api.ProblemSetResponse
	code = Post("/problem_set/create", loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    true,
		GroupId:     &groupRes.Id,
	}, &_problemSetRes)
	assert.Equal(t, code, http.StatusOK)

	var _choiceRes api.ChoiceProblemResponse
	var _choices []api.ChoiceRequest
	_choices = append(_choices, api.ChoiceRequest{
		Choice:      "A",
		Description: "test",
		IsCorrect:   true,
	})
	_choices = append(_choices, api.ChoiceRequest{
		Choice:      "B",
		Description: "test",
		IsCorrect:   false,
	})
	code = Post("/problem/choice/create", loginRes.Token, &api.ChoiceProblemCreateRequest{
		Description: "test",
		IsPublic:    false,
		Choices:     _choices,
	}, &_choiceRes)
	assert.Equal(t, code, http.StatusOK)

	var choiceRes api.ChoiceProblemResponse
	var choices []api.ChoiceRequest
	choices = append(choices, api.ChoiceRequest{
		Choice:      "A",
		Description: "test",
		IsCorrect:   true,
	})
	choices = append(choices, api.ChoiceRequest{
		Choice:      "B",
		Description: "test",
		IsCorrect:   false,
	})
	code = Post("/problem/choice/create", loginRes.Token, &api.ChoiceProblemCreateRequest{
		Description: "test",
		IsPublic:    false,
		Choices:     choices,
	}, &choiceRes)
	assert.Equal(t, code, http.StatusOK)

	var res interface{}
	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", problemSetRes.ID, choiceRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", _problemSetRes.ID, _choiceRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Put(fmt.Sprintf("/problem/choice/update"), loginRes.Token, &api.ChoiceProblemUpdateRequest{
		ID: 100,
	}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Put(fmt.Sprintf("/problem/choice/update"), _loginRes.Token, &api.ChoiceProblemUpdateRequest{
		ID: choiceRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Put(fmt.Sprintf("/problem/choice/update"), loginRes.Token, &api.ChoiceProblemUpdateRequest{
		ID:      choiceRes.ID,
		Choices: choices,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Put(fmt.Sprintf("/problem/choice/update"), _loginRes.Token, &api.ChoiceProblemUpdateRequest{
		ID: _choiceRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Put(fmt.Sprintf("/problem/choice/update"), loginRes.Token, &api.ChoiceProblemUpdateRequest{
		ID:      _choiceRes.ID,
		Choices: _choices,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/choice/answer/100"), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Get(fmt.Sprintf("/problem/choice/answer/%d", choiceRes.ID), _loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/problem/choice/answer/%d", choiceRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/choice/answer/%d", _choiceRes.ID), _loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/problem/choice/answer/%d", _choiceRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	var _blankRes api.BlankProblemResponse
	code = Post("/problem/blank/create", loginRes.Token, &api.BlankProblemCreateRequest{
		Description:   "test",
		IsPublic:      false,
		Answer:        "test",
		AnswerExplain: "test",
	}, &_blankRes)
	assert.Equal(t, code, http.StatusOK)

	var blankRes api.BlankProblemResponse
	code = Post("/problem/blank/create", loginRes.Token, &api.BlankProblemCreateRequest{
		Description:   "test",
		IsPublic:      false,
		Answer:        "test",
		AnswerExplain: "test",
	}, &blankRes)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", _problemSetRes.ID, blankRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", problemSetRes.ID, _blankRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Put(fmt.Sprintf("/problem/blank/update"), loginRes.Token, &api.BlankProblemUpdateRequest{
		ID: 100,
	}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Put(fmt.Sprintf("/problem/blank/update"), _loginRes.Token, &api.BlankProblemUpdateRequest{
		ID: blankRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Put(fmt.Sprintf("/problem/blank/update"), loginRes.Token, &api.BlankProblemUpdateRequest{
		ID: blankRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Put(fmt.Sprintf("/problem/blank/update"), _loginRes.Token, &api.BlankProblemUpdateRequest{
		ID: _blankRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Put(fmt.Sprintf("/problem/blank/update"), loginRes.Token, &api.BlankProblemUpdateRequest{
		ID: _blankRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/blank/answer/100"), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Get(fmt.Sprintf("/problem/blank/answer/%d", blankRes.ID), _loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/problem/blank/answer/%d", blankRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/blank/answer/%d", _blankRes.ID), _loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/problem/blank/answer/%d", _blankRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	var _judgeRes api.JudgeProblemResponse
	code = Post("/problem/judge/create", loginRes.Token, &api.JudgeProblemCreateRequest{
		Description: "test",
		IsPublic:    false,
		IsCorrect:   true,
	}, &_judgeRes)
	assert.Equal(t, code, http.StatusOK)

	var judgeRes api.JudgeProblemResponse
	var analysis = "test"
	code = Post("/problem/judge/create", loginRes.Token, &api.JudgeProblemCreateRequest{
		Description: "test",
		IsPublic:    false,
		IsCorrect:   true,
		Analysis:    &analysis,
	}, &judgeRes)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", _problemSetRes.ID, judgeRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", problemSetRes.ID, _judgeRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Put(fmt.Sprintf("/problem/judge/update"), loginRes.Token, &api.JudgeProblemUpdateRequest{
		ID: 100,
	}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Put(fmt.Sprintf("/problem/judge/update"), _loginRes.Token, &api.JudgeProblemUpdateRequest{
		ID: judgeRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Put(fmt.Sprintf("/problem/judge/update"), loginRes.Token, &api.JudgeProblemUpdateRequest{
		ID: judgeRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Put(fmt.Sprintf("/problem/judge/update"), _loginRes.Token, &api.JudgeProblemUpdateRequest{
		ID: _judgeRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Put(fmt.Sprintf("/problem/judge/update"), loginRes.Token, &api.JudgeProblemUpdateRequest{
		ID: _judgeRes.ID,
	}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/judge/answer/100"), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Get(fmt.Sprintf("/problem/judge/answer/%d", judgeRes.ID), _loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/problem/judge/answer/%d", judgeRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem/judge/answer/%d", _judgeRes.ID), _loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Get(fmt.Sprintf("/problem/judge/answer/%d", _judgeRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	var problemSetR api.ProblemSetResponse
	code = Post("/problem_set/create", _loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    true,
	}, &problemSetR)
	assert.Equal(t, code, http.StatusOK)

	var _problemSetR api.ProblemSetResponse
	code = Post("/problem_set/create", _loginRes.Token, &api.ProblemSetCreateRequest{
		Name:        "test",
		Description: "test",
		IsPublic:    true,
		GroupId:     &_groupRes.Id,
	}, &_problemSetR)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/migrate/%d?problem_id=%d", 100, 100), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/problem_set/migrate/%d?problem_id=%d", problemSetRes.ID, 100), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/problem_set/migrate/%d?problem_id=%d", _problemSetRes.ID, 100), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/problem_set/migrate/%d?problem_id=%d", problemSetRes.ID, 100), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Post(fmt.Sprintf("/problem_set/migrate/%d?problem_id=%d", _problemSetRes.ID, choiceRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/migrate/%d?problem_id=%d", _problemSetR.ID, choiceRes.ID), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/problem_set/migrate/%d?problem_id=%d", problemSetRes.ID, blankRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/migrate/%d?problem_id=%d", problemSetR.ID, blankRes.ID), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Post(fmt.Sprintf("/problem_set/migrate/%d?problem_id=%d", problemSetRes.ID, judgeRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/migrate/%d?problem_id=%d", problemSetR.ID, judgeRes.ID), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/problem_set/remove/%d?problem_id=%d", 100, 100), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/problem_set/remove/%d?problem_id=%d", problemSetRes.ID, 100), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/problem_set/remove/%d?problem_id=%d", _problemSetRes.ID, 100), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/problem_set/remove/%d?problem_id=%d", problemSetRes.ID, 100), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/problem_set/remove/%d?problem_id=%d", problemSetRes.ID, choiceRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", problemSetRes.ID, choiceRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/problem/choice/delete/100"), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/problem/choice/delete/%d", choiceRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/problem/choice/delete/%d", _choiceRes.ID), _loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/problem/choice/delete/%d", _choiceRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/problem/blank/delete/%d", _blankRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/problem/judge/delete/%d", _judgeRes.ID), loginRes.Token, struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Delete(fmt.Sprintf("/problem_set/delete/%d", 100), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusNotFound)

	code = Delete(fmt.Sprintf("/problem_set/delete/%d", problemSetRes.ID), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/problem_set/delete/%d", _problemSetRes.ID), _loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusForbidden)

	code = Delete(fmt.Sprintf("/problem_set/delete/%d", problemSetRes.ID), loginRes.Token, &struct{}{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem_set/statistic/wrong_count?id=%d", problemSetRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)

	code = Get(fmt.Sprintf("/problem_set/statistic/fav_count?id=%d", problemSetRes.ID), loginRes.Token, map[string][]string{}, &res)
	assert.Equal(t, code, http.StatusOK)
}
