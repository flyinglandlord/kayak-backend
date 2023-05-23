package model

import "time"

type ProblemSet struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UserId      int       `json:"user_id" db:"user_id"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	GroupId     int       `json:"group_id" db:"group_id"`
	AreaId      int       `json:"area_id" db:"area_id"`
	LikeCount   int       `json:"like_count" db:"like_count"`
}
