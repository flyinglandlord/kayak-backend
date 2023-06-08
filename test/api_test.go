package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"kayak-backend/api"
	"kayak-backend/global"
	"kayak-backend/model"
	"kayak-backend/utils"
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"
)

var stages = [][]func(*testing.T){
	{testPing, testLogin, testRegister, testChangePassword, testResetPassword, testLogout},
	{testGetUserInfoById, testGetUserInfo, testUpdateUserInfo, testGetUserWrongRecords},
	{testGetWrongProblemSet, testGetFavoriteProblemSet},
	{testGetFeaturedProblemSet, testGetFeaturedNote, testGetFeaturedGroup},
	{testSearchProblemSets, testSearchGroups, testSearchNotes},
	{testGetProblemSets, testCreateProblemSet, testUpdateProblemSet},
	{testGetProblemsInProblemSet, testAddProblemToProblemSet},
	{testMigrateProblemToProblemSet},
	{testGetChoiceProblems, testGetBlankProblems, testGetJudgeProblems},
	{testCreateWrongRecord},
	{testGetNotes, testCreateNote},
	{testAddNoteReview},
	{testGetGroups, testCreateGroup},
	{testCreateDiscussion},
	{testAddDiscussionReview},
}

func goTestWithWait(wg *sync.WaitGroup, t *testing.T, f func(t *testing.T)) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("--> run " + runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		f(t)
	}()
}

var tokens []string

func Benchmark_TimeConsumingFunction(b *testing.B) {
	b.SetParallelism(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			epoch()
		}
	})
}

func TestMain(m *testing.M) {
	readConfig()
	global.Redis.FlushDB(context.Background())
	fileContent, err := ioutil.ReadFile("../init.sql")
	if err != nil {
		panic(err)
	}
	if _, err := global.Database.Exec(string(fileContent)); err != nil {
		panic(err)
	}
	tx := global.Database.MustBegin()
	if err != nil {
		panic(err)
	}

	sqlString := `INSERT INTO "user" (name, created_at, password, nick_name, email) VALUES ($1, now(), $2, $3, $4)`
	user := model.User{
		ID:       1,
		Name:     fmt.Sprintf("user0"),
		NickName: fmt.Sprintf("user0"),
		Email:    fmt.Sprintf("user%d@boat4study.com", 0),
		Password: fmt.Sprintf("%d-pwd", 0),
	}
	randomUsers = append(randomUsers, user)
	encryptPassword, _ := utils.EncryptPassword(user.Password)
	if _, err := tx.Exec(sqlString, user.Name, encryptPassword, user.NickName, user.Email); err != nil {
		panic(err)
	}
	sqlString = `INSERT INTO "group" (name, description, invitation, created_at, user_id, area_id) VALUES ($1, $2, $3, now(), $4, $5)`
	if _, err := tx.Exec(sqlString, "test", "test", 0, 1, 1); err != nil {
		panic(err)
	}
	_sqlString := `INSERT INTO group_member (group_id, user_id, created_at, is_admin, is_owner) VALUES ($1, $2, now(), $3, $4)`
	if _, err := tx.Exec(_sqlString, 1, 1, false, true); err != nil {
		panic(err)
	}
	sqlString = `INSERT INTO "user" (name, created_at, password, nick_name, email) VALUES ($1, now(), $2, $3, $4)`
	for i := 1; i < 100; i++ {
		user := model.User{
			ID:       i + 1,
			Name:     fmt.Sprintf("user%d", i),
			NickName: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@boat4study.com", i),
			Password: fmt.Sprintf("%d-pwd", i),
		}
		randomUsers = append(randomUsers, user)
		encryptPassword, _ := utils.EncryptPassword(user.Password)
		if _, err := tx.Exec(sqlString, user.Name, encryptPassword, user.NickName, user.Email); err != nil {
			panic(err)
		}
		_sqlString := `INSERT INTO group_member (group_id, user_id, created_at, is_admin, is_owner) VALUES ($1, $2, now(), $3, $4)`
		if _, err := tx.Exec(_sqlString, 1, i+1, false, false); err != nil {
			panic(err)
		}
	}
	sqlString = `INSERT INTO problem_set (name, description, created_at, updated_at, 
        user_id, is_public, area_id, group_id) VALUES ($1, $2, now(), now(), $3, $4, $5, $6)`
	if _, err := tx.Exec(sqlString, "test", "test", 1, false, 1, 1); err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		loginRes := api.LoginResponse{}
		user := randomUsers[i]
		Post("/login", "", &api.LoginInfo{
			UserName: user.Name,
			Password: user.Password,
		}, &loginRes)
		tokens = append(tokens, loginRes.Token)
	}
	exit := m.Run()
	os.Exit(exit)
}

/*
func TestBasicFunction(t *testing.T) {
	wg := sync.WaitGroup{}
	for _, group := range stages {
		InitDatabase()
		for _, f := range group {
			goTestWithWait(&wg, t, f)
		}
		wg.Wait()
	}
}*/

func epoch() {
	user := randomUser()
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
	Post("/problem/choice/create", tokens[user.ID-1], &api.ChoiceProblemCreateRequest{
		Description: "test",
		IsPublic:    false,
		Choices:     _choices,
	}, &_choiceRes)
	var res interface{}
	Post(fmt.Sprintf("/problem_set/add/%d?problem_id=%d", 1, _choiceRes.ID), tokens[user.ID-1], &struct{}{}, &res)
}
