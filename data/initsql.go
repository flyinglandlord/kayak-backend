package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"reflect"
	"time"
)

func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

var DB *sqlx.DB

func InitSql(Addr string, Port int, User string, Password string, Database string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Addr,
		Port,
		User,
		Password,
		Database)

	fmt.Println(psqlInfo)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	DB = db
}

type DanxuanQuestion struct {
	Description string `json:"description"`
	Answer      int    `json:"answer"`
	Options     struct {
		A string `json:"0"`
		B string `json:"1"`
		C string `json:"2"`
		D string `json:"3"`
	} `json:"options"`
}

type PanduanQuestion struct {
	Description string `json:"description"`
	Answer      string `json:"answer"`
	// Type        string `json:"type"`
}

type DuoxuanQuestion struct {
	Description string `json:"description"`
	Answer      []int  `json:"answer"`
	Options     struct {
		A string `json:"0"`
		B string `json:"1"`
		C string `json:"2"`
		D string `json:"3"`
	} `json:"options"`
}

func loadDanxuanQuestion(problemSetId int) {
	// 读取单选题
	jsonFile, err := os.Open("computer_basic_danxuan.json")
	if err != nil {
		panic(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			panic(err)
		}
	}(jsonFile)

	// 解析json文件
	var danxuanQuestions []DanxuanQuestion
	err = json.NewDecoder(jsonFile).Decode(&danxuanQuestions)
	if err != nil {
		panic(err)
	}

	// 插入数据库
	for _, danxuanQuestion := range danxuanQuestions {
		sqlString := "INSERT INTO problem_type (description, created_at, updated_at, user_id, is_public, problem_type_id, analysis) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
		var problemId int
		if err := DB.Get(&problemId, sqlString, danxuanQuestion.Description, time.Now(), time.Now(), 0, true, 0, ""); err != nil {
			panic(err)
		}
		sqlString = "INSERT INTO problem_choice (id, choice, description, is_correct) VALUES ($1, $2, $3, $4)"
		if _, err := DB.Exec(sqlString, problemId, "A", danxuanQuestion.Options.A, danxuanQuestion.Answer == 0); err != nil {
			panic(err)
		}
		if _, err := DB.Exec(sqlString, problemId, "B", danxuanQuestion.Options.B, danxuanQuestion.Answer == 1); err != nil {
			panic(err)
		}
		if _, err := DB.Exec(sqlString, problemId, "C", danxuanQuestion.Options.C, danxuanQuestion.Answer == 2); err != nil {
			panic(err)
		}
		if _, err := DB.Exec(sqlString, problemId, "D", danxuanQuestion.Options.D, danxuanQuestion.Answer == 3); err != nil {
			panic(err)
		}
		sqlString = "INSERT INTO problem_in_problem_set (problem_set_id, problem_id) VALUES ($1, $2)"
		if _, err := DB.Exec(sqlString, problemSetId, problemId); err != nil {
			panic(err)
		}
	}

	fmt.Println("Finish loading danxuan_question.json")
}

func loadPanduanQuestion(problemSetId int) {
	jsonFile, err := os.Open("computer_basic_panduan.json")
	if err != nil {
		panic(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			panic(err)
		}
	}(jsonFile)

	var panduanQuestions []PanduanQuestion
	err = json.NewDecoder(jsonFile).Decode(&panduanQuestions)
	if err != nil {
		panic(err)
	}

	for _, panduanQuestion := range panduanQuestions {
		sqlString := "INSERT INTO problem_type (description, created_at, updated_at, user_id, is_public, problem_type_id, analysis) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
		var problemId int
		if err := DB.Get(&problemId, sqlString, panduanQuestion.Description, time.Now(), time.Now(), 0, true, 2, ""); err != nil {
			panic(err)
		}
		sqlString = "INSERT INTO problem_judge (id, is_correct) VALUES ($1, $2)"
		if _, err := DB.Exec(sqlString, problemId, panduanQuestion.Answer == "正确"); err != nil {
			panic(err)
		}
		sqlString = "INSERT INTO problem_in_problem_set (problem_set_id, problem_id) VALUES ($1, $2)"
		if _, err := DB.Exec(sqlString, problemSetId, problemId); err != nil {
			panic(err)
		}
	}

	fmt.Println("Finish loading panduan_question.json")
}

func loadDuoxuanQuestion(problemSetId int) {
	jsonFile, err := os.Open("computer_basic_duoxuan.json")
	if err != nil {
		panic(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			panic(err)
		}
	}(jsonFile)

	var duoxuanQuestions []DuoxuanQuestion
	err = json.NewDecoder(jsonFile).Decode(&duoxuanQuestions)
	if err != nil {
		panic(err)
	}

	for _, duoxuanQuestion := range duoxuanQuestions {
		sqlString := `INSERT INTO problem_type (description, created_at, updated_at, user_id, is_public, problem_type_id, analysis) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
		var problemId int
		if err := DB.Get(&problemId, sqlString, duoxuanQuestion.Description, time.Now(), time.Now(), 0, true, 0, ""); err != nil {
			panic(err)
		}
		sqlString = `INSERT INTO problem_choice (id, choice, description, is_correct) VALUES ($1, $2, $3, $4)`
		if _, err := DB.Exec(sqlString, problemId, "A", duoxuanQuestion.Options.A, Contain(duoxuanQuestion.Answer, 0)); err != nil {
			panic(err)
		}
		if _, err := DB.Exec(sqlString, problemId, "B", duoxuanQuestion.Options.B, Contain(duoxuanQuestion.Answer, 1)); err != nil {
			panic(err)
		}
		if _, err := DB.Exec(sqlString, problemId, "C", duoxuanQuestion.Options.C, Contain(duoxuanQuestion.Answer, 2)); err != nil {
			panic(err)
		}
		if _, err := DB.Exec(sqlString, problemId, "D", duoxuanQuestion.Options.D, Contain(duoxuanQuestion.Answer, 3)); err != nil {
			panic(err)
		}
		sqlString = "INSERT INTO problem_in_problem_set (problem_set_id, problem_id) VALUES ($1, $2)"
		if _, err := DB.Exec(sqlString, problemSetId, problemId); err != nil {
			panic(err)
		}
	}

	fmt.Println("Finish loading duoxuan_question.json")
}

func main() {
	// Init Database
	// **************************************************************(secret)
	sqlString := `INSERT INTO "user" (id, name, password, created_at, nick_name) VALUES ($1, $2, $3, $4, $5)`
	/*if _, err := DB.Exec(sqlString, 0, "problem_set_maker", "123456", time.Now(), "problem_set_maker"); err != nil {
		panic(err)
	}*/

	sqlString = "INSERT INTO problem_set (name, description, created_at, updated_at, user_id, is_public, group_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	var problemSetId int
	if err := DB.Get(&problemSetId, sqlString, "计算机基础", "计算机基础题库", time.Now(), time.Now(), 0, true, 0); err != nil {
		panic(err)
	}
	fmt.Println(problemSetId)
	loadDanxuanQuestion(problemSetId)
	loadDuoxuanQuestion(problemSetId)
	loadPanduanQuestion(problemSetId)
}
