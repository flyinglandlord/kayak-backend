package model

import "time"

type GroupMember struct {
	GroupId   int       `json:"group_id" db:"group_id"`
	UserId    int       `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	IsOwner   bool      `json:"is_owner" db:"is_owner"`
}
