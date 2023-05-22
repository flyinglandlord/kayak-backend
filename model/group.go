package model

import "time"

type Group struct {
	Id          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Invitation  string    `json:"invitation" db:"invitation"`
	UserId      int       `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	AreaId      int       `json:"area_id" db:"area_id"`
}
