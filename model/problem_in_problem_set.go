package model

type ProblemInProblemSet struct {
	ProblemSetId int `json:"problem_set_id" db:"problem_set_id"`
	ProblemId    int `json:"problem_id" db:"problem_id"`
}
