package model

import (
	"time"
)

const (
	SessionTokenSize = 32
)

type Session struct {
	Id        string     `db:"id" json:"id"`
	UserId    *string    `db:"user_id" json:"userId,omitempty"`
	Token     *string    `db:"token" json:"token"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt,omitempty"`
}

func DefaultSession() *Session {
	t := time.Now()
	token := generateRandomString(SessionTokenSize)
	return &Session{
		generateID(),
		nil,
		&token,
		&t,
	}
}
