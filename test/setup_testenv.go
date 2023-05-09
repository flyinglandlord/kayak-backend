package test

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"io/ioutil"
	"kayak-backend/api"
	"kayak-backend/global"
	"kayak-backend/model"
	"kayak-backend/utils"
	"time"
)

var initUser = []model.User{
	{Name: "test1", CreatedAt: time.Now(), NickName: "test1", Email: "test1@boat4study.com"},
	{Name: "test2", CreatedAt: time.Now(), NickName: "test2", Email: "test2@boat4study.com"},
	{Name: "test3", CreatedAt: time.Now(), NickName: "test3", Email: "test3@boat4study.com"},
	{Name: "test4", CreatedAt: time.Now(), NickName: "test4", Email: "test4@boat4study.com"},
	{Name: "test5", CreatedAt: time.Now(), NickName: "test5", Email: "test5@boat4study.com"},
	{Name: "test6", CreatedAt: time.Now(), NickName: "test6", Email: "test6@boat4study.com"},
	{Name: "test7", CreatedAt: time.Now(), NickName: "test7", Email: "test7@boat4study.com"},
	{Name: "test8", CreatedAt: time.Now(), NickName: "test8", Email: "test8@boat4study.com"},
}

var initProblemInProblemSet = []model.ProblemInProblemSet{
	{ProblemSetId: 1, ProblemId: 1},
	{ProblemSetId: 2, ProblemId: 2},
	{ProblemSetId: 3, ProblemId: 3},
	{ProblemSetId: 4, ProblemId: 4},
	{ProblemSetId: 5, ProblemId: 5},
	{ProblemSetId: 6, ProblemId: 6},
	{ProblemSetId: 7, ProblemId: 7},
	{ProblemSetId: 8, ProblemId: 8},
}

var initProblemType = []model.ProblemType{
	{ID: 1, Description: "problem1", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 1, ProblemTypeId: 0, IsPublic: true},
	{ID: 2, Description: "problem2", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 2, ProblemTypeId: 1, IsPublic: true},
	{ID: 3, Description: "problem3", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 3, ProblemTypeId: 0, IsPublic: true},
	{ID: 4, Description: "problem4", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 4, ProblemTypeId: 1, IsPublic: true},
	{ID: 5, Description: "problem5", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 5, ProblemTypeId: 0, IsPublic: false},
	{ID: 6, Description: "problem6", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 6, ProblemTypeId: 1, IsPublic: false},
	{ID: 7, Description: "problem7", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 7, ProblemTypeId: 0, IsPublic: false},
	{ID: 8, Description: "problem8", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 8, ProblemTypeId: 1, IsPublic: false},
}

var initProblemChoice = []model.ProblemChoice{
	{ID: 1, Choice: "A", Description: "description1", IsCorrect: true},
	{ID: 1, Choice: "B", Description: "description2", IsCorrect: false},
	{ID: 1, Choice: "C", Description: "description3", IsCorrect: false},
	{ID: 1, Choice: "D", Description: "description4", IsCorrect: false},
	{ID: 3, Choice: "A", Description: "description1", IsCorrect: false},
	{ID: 3, Choice: "B", Description: "description2", IsCorrect: true},
	{ID: 3, Choice: "C", Description: "description3", IsCorrect: false},
	{ID: 3, Choice: "D", Description: "description4", IsCorrect: false},
	{ID: 5, Choice: "A", Description: "description1", IsCorrect: false},
	{ID: 5, Choice: "B", Description: "description2", IsCorrect: false},
	{ID: 5, Choice: "C", Description: "description3", IsCorrect: true},
	{ID: 5, Choice: "D", Description: "description4", IsCorrect: false},
	{ID: 7, Choice: "A", Description: "description1", IsCorrect: false},
	{ID: 7, Choice: "B", Description: "description2", IsCorrect: false},
	{ID: 7, Choice: "C", Description: "description3", IsCorrect: false},
	{ID: 7, Choice: "D", Description: "description4", IsCorrect: true},
}

var initProblemAnswer = []model.ProblemAnswer{
	{ID: 2, Answer: "problem2_answer"},
	{ID: 4, Answer: "problem4_answer"},
	{ID: 6, Answer: "problem6_answer"},
	{ID: 8, Answer: "problem8_answer"},
}

var initProblemSet = []model.ProblemSet{
	{ID: 1, Name: "name1", Description: "description1", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 1, IsPublic: false},
	{ID: 2, Name: "name2", Description: "description2", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 2, IsPublic: false},
	{ID: 3, Name: "name3", Description: "description3", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 3, IsPublic: false},
	{ID: 4, Name: "name4", Description: "description4", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 4, IsPublic: false},
	{ID: 5, Name: "name5", Description: "description5", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 5, IsPublic: true},
	{ID: 6, Name: "name6", Description: "description6", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 6, IsPublic: true},
	{ID: 7, Name: "name7", Description: "description7", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 7, IsPublic: true},
	{ID: 8, Name: "name8", Description: "description8", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 8, IsPublic: true},
}

var initNote = []model.Note{
	{ID: 1, Title: "title1", Content: "content1", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 1, IsPublic: false},
	{ID: 2, Title: "title2", Content: "content2", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 2, IsPublic: false},
	{ID: 3, Title: "title3", Content: "content3", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 3, IsPublic: false},
	{ID: 4, Title: "title4", Content: "content4", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 4, IsPublic: false},
	{ID: 5, Title: "title5", Content: "content5", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 5, IsPublic: true},
	{ID: 6, Title: "title6", Content: "content6", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 6, IsPublic: true},
	{ID: 7, Title: "title7", Content: "content7", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 7, IsPublic: true},
	{ID: 8, Title: "title8", Content: "content8", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 8, IsPublic: true},
}

var initNoteReview = []model.NoteReview{
	{ID: 1, NoteId: 1, UserId: 1, Title: "title1", Content: "content1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, NoteId: 1, UserId: 2, Title: "title2", Content: "content2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 3, NoteId: 2, UserId: 1, Title: "title3", Content: "content3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 4, NoteId: 2, UserId: 3, Title: "title4", Content: "content4", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 5, NoteId: 3, UserId: 4, Title: "title5", Content: "content5", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 6, NoteId: 3, UserId: 2, Title: "title6", Content: "content6", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 7, NoteId: 4, UserId: 5, Title: "title7", Content: "content7", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 8, NoteId: 4, UserId: 6, Title: "title8", Content: "content8", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 9, NoteId: 5, UserId: 3, Title: "title9", Content: "content9", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 10, NoteId: 5, UserId: 7, Title: "title10", Content: "content10", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 11, NoteId: 6, UserId: 8, Title: "title11", Content: "content11", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 12, NoteId: 6, UserId: 4, Title: "title12", Content: "content12", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 13, NoteId: 7, UserId: 5, Title: "title13", Content: "content13", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 14, NoteId: 7, UserId: 6, Title: "title14", Content: "content14", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 15, NoteId: 8, UserId: 7, Title: "title15", Content: "content15", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 16, NoteId: 8, UserId: 8, Title: "title16", Content: "content16", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

var initWrongRecord = []model.WrongRecord{
	{ProblemId: 1, UserId: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 1},
	{ProblemId: 2, UserId: 2, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 2},
	{ProblemId: 3, UserId: 3, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 3},
	{ProblemId: 4, UserId: 4, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 4},
	{ProblemId: 5, UserId: 5, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 5},
	{ProblemId: 6, UserId: 6, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 6},
	{ProblemId: 7, UserId: 7, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 7},
	{ProblemId: 8, UserId: 8, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 8},
	{ProblemId: 1, UserId: 5, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 9},
	{ProblemId: 2, UserId: 6, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 10},
	{ProblemId: 3, UserId: 7, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 11},
	{ProblemId: 4, UserId: 8, CreatedAt: time.Now(), UpdatedAt: time.Now(), Count: 12},
}

func readConfig() {
	viper.SetConfigName("config_test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("PostgresHost"),
		viper.GetInt("PostgresPort"),
		viper.GetString("PostgresUsername"),
		viper.GetString("PostgresPassword"),
		viper.GetString("PostgresDatabase"))
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	global.Database = db

	global.Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("RedisHost"), viper.GetInt("RedisPort")),
		Password: viper.GetString("RedisPassword"),
		DB:       0,
	})
	status := global.Redis.Ping(context.Background())
	if status.Err() != nil {
		panic(status.Err())
	}

	global.Router = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, global.TokenHeader)
	global.Router.Use(cors.New(corsConfig))
	global.Router.Use(global.Authenticate)
	api.InitRoute()
}

func InitUserTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO "user" (name, created_at, password, nick_name, email) VALUES ($1, now(), $2, $3, $4)`
	for i := range initUser {
		initUser[i].Password = fmt.Sprintf("%s-pwd", initUser[i].Name)
		encryptPassword, _ := utils.EncryptPassword(initUser[i].Password)
		if _, err := tx.Exec(sqlString, initUser[i].Name, encryptPassword, initUser[i].NickName, initUser[i].Email); err != nil {
			return err
		}
	}
	return nil
}

func InitProblemInProblemSetTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO problem_in_problem_set (problem_id, problem_set_id) VALUES ($1, $2)`
	for i := range initProblemInProblemSet {
		if _, err := tx.Exec(sqlString, initProblemInProblemSet[i].ProblemId, initProblemInProblemSet[i].ProblemSetId); err != nil {
			return err
		}
	}
	return nil
}

func InitProblemTypeTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO problem_type (description, created_at, updated_at, user_id, problem_type_id, is_public) VALUES ($1, now(), now(), $2, $3, $4)`
	for i := range initProblemType {
		if _, err := tx.Exec(sqlString, initProblemType[i].Description, initProblemType[i].UserId,
			initProblemType[i].ProblemTypeId, initProblemType[i].IsPublic); err != nil {
			return err
		}
	}
	return nil
}

func InitProblemChoiceTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO problem_choice (id, choice, description, is_correct) VALUES ($1, $2, $3, $4)`
	for i := range initProblemChoice {
		if _, err := tx.Exec(sqlString, initProblemChoice[i].ID, initProblemChoice[i].Choice, initProblemChoice[i].Description, initProblemChoice[i].IsCorrect); err != nil {
			return err
		}
	}
	return nil
}

func InitProblemAnswerTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO problem_answer (id, answer) VALUES ($1, $2)`
	for i := range initProblemAnswer {
		if _, err := tx.Exec(sqlString, initProblemAnswer[i].ID, initProblemAnswer[i].Answer); err != nil {
			return err
		}
	}
	return nil
}

func InitProblemSetTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO problem_set (name, description, created_at, updated_at, user_id, is_public) VALUES ($1, $2, $3, $4, $5, $6)`
	for i := range initProblemSet {
		if _, err := tx.Exec(sqlString, initProblemSet[i].Name,
			initProblemSet[i].Description, initProblemSet[i].CreatedAt, initProblemSet[i].UpdatedAt,
			initProblemSet[i].UserId, initProblemSet[i].IsPublic); err != nil {
			return err
		}
	}
	return nil
}

func InitNoteTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO note (title, content, created_at, updated_at, user_id, is_public) VALUES ($1, $2, now(), now(), $3, $4)`
	for i := range initNote {
		if _, err := tx.Exec(sqlString, initNote[i].Title, initNote[i].Content, initNote[i].UserId, initNote[i].IsPublic); err != nil {
			return err
		}
	}
	return nil
}

func InitNoteReviewTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO note_review (title, content, created_at, updated_at, user_id, note_id) VALUES ($1, $2, now(), now(), $3, $4)`
	for i := range initNoteReview {
		if _, err := tx.Exec(sqlString, initNoteReview[i].Title, initNoteReview[i].Content, initNoteReview[i].UserId, initNoteReview[i].NoteId); err != nil {
			return err
		}
	}
	return nil
}

func InitWrongRecord(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO user_wrong_record (user_id, problem_id, created_at, updated_at, count) VALUES ($1, $2, $3, $4, $5)`
	for i := range initWrongRecord {
		if _, err := tx.Exec(sqlString, initWrongRecord[i].UserId, initWrongRecord[i].ProblemId,
			initWrongRecord[i].CreatedAt, initWrongRecord[i].UpdatedAt, initWrongRecord[i].Count); err != nil {
			return err
		}
	}
	return nil

}

var initFuncList = []func(tx *sqlx.Tx) error{
	InitUserTable,
	InitProblemTypeTable,
	InitProblemChoiceTable,
	InitProblemAnswerTable,
	InitProblemSetTable,
	InitProblemTypeTable,
	InitProblemInProblemSetTable,
	InitNoteTable,
	InitNoteReviewTable,
	InitWrongRecord,
}

func InitDatabase() {
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
	for _, f := range initFuncList {
		// fmt.Println("Finish init table: ", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		if err := f(tx); err != nil {
			panic(err)
		}
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
