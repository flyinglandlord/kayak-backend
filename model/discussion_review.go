package model

type DiscussionReview struct {
	ID           int    `json:"id" db:"id"`
	DiscussionId int    `json:"discussion_id" db:"discussion_id"`
	UserId       int    `json:"user_id" db:"user_id"`
	Title        string `json:"title" db:"title"`
	Content      string `json:"content" db:"content"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	UpdatedAt    string `json:"updated_at" db:"updated_at"`
}
