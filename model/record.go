package model

import "time"

type WrongRecord struct {
	ProblemId int       `json:"problem_id" db:"problem_id"`
	UserId    int       `json:"user_id" db:"user_id"`
	Count     int       `json:"count" db:"count"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
