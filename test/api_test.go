package test

import (
	"fmt"
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

func TestMain(m *testing.M) {
	readConfig()
	exit := m.Run()
	os.Exit(exit)
}

func TestBasicFunction(t *testing.T) {
	wg := sync.WaitGroup{}
	for _, group := range stages {
		InitDatabase()
		for _, f := range group {
			goTestWithWait(&wg, t, f)
		}
		wg.Wait()
	}
}
