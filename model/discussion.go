package model

type Discussion struct {
	ID        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	UserId    int    `json:"user_id" db:"user_id"`
	GroupId   int    `json:"group_id" db:"group_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
	IsPublic  bool   `json:"is_public" db:"is_public"`
	LikeCount int    `json:"like_count" db:"like_count"`
}
