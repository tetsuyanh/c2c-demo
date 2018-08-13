package model

import "time"

type (
	User struct {
		ID        string     `db:"id" json:"id"`
		CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	}
)

func NewUser() *User {
	t := time.Now()
	return &User{
		generateID(),
		&t,
	}
}
