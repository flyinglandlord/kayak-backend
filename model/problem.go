package model

import "time"

type ProblemType struct {
	ID            int       `json:"id" db:"id"`
	Description   string    `json:"description" db:"description"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	UserId        int       `json:"user_id" db:"user_id"`
	ProblemTypeId int       `json:"problem_type_id" db:"problem_type_id"`
	IsPublic      bool      `json:"is_public" db:"is_public"`
	Analysis      *string   `json:"analysis" db:"analysis"`
}

type ProblemChoice struct {
	ID          int    `json:"id" db:"id"`
	Choice      string `json:"choice" db:"choice"`
	Description string `json:"description" db:"description"`
	IsCorrect   bool   `json:"is_correct" db:"is_correct"`
}

type ProblemAnswer struct {
	ID     int    `json:"id" db:"id"`
	Answer string `json:"answer" db:"answer"`
}

type ProblemJudge struct {
	ID        int  `json:"id" db:"id"`
	IsCorrect bool `json:"is_correct" db:"is_correct"`
}
