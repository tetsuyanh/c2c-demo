package model

import (
	"time"
)

type (
	User struct {
		Id        string     `db:"id" json:"id"`
		CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	}
)

func DefaultUser() *User {
	t := time.Now()
	return &User{
		generateID(),
		&t,
	}
}
