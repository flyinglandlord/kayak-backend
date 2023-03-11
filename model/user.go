package model

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	UnionID   *string   `json:"union_id" db:"union_id"`
	Name      string    `json:"name" db:"name"`
	Email     *string   `json:"email" db:"email"`
	Phone     *string   `json:"phone" db:"phone"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	AvatarURL string    `json:"avatar_url" db:"avatar_url"`
}
