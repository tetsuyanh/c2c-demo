package model

import (
	"time"
)

const (
	SessionTokenSize = 32
)

type Session struct {
	Id        string    `db:"id" json:"id"`
	UserId    string    `db:"user_id" json:"userId,omitempty"`
	Token     string    `db:"token" json:"token"`
	CreatedAt time.Time `db:"created_at" json:"createdAt,omitempty"`
}

func DefaultSession() *Session {
	return &Session{
		generateID(),
		"",
		generateRandomString(SessionTokenSize),
		time.Now(),
	}
}
