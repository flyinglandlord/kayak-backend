package model

type GroupApplication struct {
	ID        int     `json:"id" db:"id"`
	UserId    int     `json:"user_id" db:"user_id"`
	GroupId   int     `json:"group_id" db:"group_id"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	Status    int     `json:"status" db:"status"` // 0: pending, 1: accepted, 2: rejected
	Message   *string `json:"message" db:"message"`
}
