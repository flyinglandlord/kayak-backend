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
	{ID: 1, Name: "test1", CreatedAt: time.Now()},
	{ID: 2, Name: "test2", CreatedAt: time.Now()},
	{ID: 3, Name: "test3", CreatedAt: time.Now()},
	{ID: 4, Name: "test4", CreatedAt: time.Now()},
	{ID: 5, Name: "test5", CreatedAt: time.Now()},
}

var initProblemType = []model.ProblemType{
	{ID: 1, Description: "problem1", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: 1, ProblemTypeID: 0, IsPublic: false},
	{ID: 2, Description: "problem2", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: 2, ProblemTypeID: 1, IsPublic: false},
	{ID: 3, Description: "problem3", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: 3, ProblemTypeID: 0, IsPublic: false},
	{ID: 4, Description: "problem4", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: 4, ProblemTypeID: 1, IsPublic: false},
	{ID: 5, Description: "problem5", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: 5, ProblemTypeID: 0, IsPublic: false},
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
}

var initProblemAnswer = []model.ProblemAnswer{
	{ID: 2, Answer: "problem2_answer"},
	{ID: 4, Answer: "problem4_answer"},
}

var initNote = []model.Note{
	{ID: 1, Title: "title1", Content: "content1", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 1, IsPublic: false},
	{ID: 2, Title: "title2", Content: "content2", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 2, IsPublic: false},
	{ID: 3, Title: "title3", Content: "content3", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 3, IsPublic: false},
	{ID: 4, Title: "title4", Content: "content4", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 4, IsPublic: false},
	{ID: 5, Title: "title5", Content: "content5", CreatedAt: time.Now(), UpdatedAt: time.Now(), UserId: 5, IsPublic: false},
}

var initNoteReview = []model.NoteReview{
	{ID: 1, NoteId: 1, UserId: 2, Title: "title1", Content: "content1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, NoteId: 1, UserId: 3, Title: "title2", Content: "content2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 3, NoteId: 2, UserId: 1, Title: "title2", Content: "content2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 4, NoteId: 2, UserId: 4, Title: "title3", Content: "content3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 5, NoteId: 3, UserId: 2, Title: "title3", Content: "content3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 6, NoteId: 3, UserId: 5, Title: "title4", Content: "content4", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 7, NoteId: 4, UserId: 5, Title: "title4", Content: "content4", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 8, NoteId: 4, UserId: 1, Title: "title5", Content: "content5", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 9, NoteId: 5, UserId: 3, Title: "title5", Content: "content5", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 10, NoteId: 5, UserId: 2, Title: "title5", Content: "content5", CreatedAt: time.Now(), UpdatedAt: time.Now()},
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
	sqlString := `INSERT INTO "user" (id, name, created_at, password) VALUES ($1, $2, now(), $3)`
	for i := range initUser {
		initUser[i].Password = fmt.Sprintf("%s-pwd", initUser[i].Name)
		encryptPassword, _ := utils.EncryptPassword(initUser[i].Password)
		if _, err := tx.Exec(sqlString, initUser[i].ID, initUser[i].Name, encryptPassword); err != nil {
			return err
		}
	}
	return nil
}

func InitProblemTypeTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO problem_type (id, description, created_at, updated_at, user_id, problem_type_id, is_public) VALUES ($1, $2, now(), now(), $3, $4, $5)`
	for i := range initProblemType {
		if _, err := tx.Exec(sqlString, initProblemType[i].ID, initProblemType[i].Description, initProblemType[i].UserID,
			initProblemType[i].ProblemTypeID, initProblemType[i].IsPublic); err != nil {
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

func InitNoteTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO note (id, title, content, created_at, updated_at, user_id, is_public) VALUES ($1, $2, $3, now(), now(), $4, $5)`
	for i := range initNote {
		if _, err := tx.Exec(sqlString, initNote[i].ID, initNote[i].Title, initNote[i].Content, initNote[i].UserId, initNote[i].IsPublic); err != nil {
			return err
		}
	}
	return nil
}

func InitNoteReviewTable(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO note_review (id, title, content, created_at, updated_at, user_id, note_id) VALUES ($1, $2, $3, now(), now(), $4, $5)`
	for i := range initNoteReview {
		if _, err := tx.Exec(sqlString, initNoteReview[i].ID, initNoteReview[i].Title, initNoteReview[i].Content, initNoteReview[i].UserId, initNoteReview[i].NoteId); err != nil {
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
	InitNoteTable,
	InitNoteReviewTable,
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
