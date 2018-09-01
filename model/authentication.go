package model

import (
	"time"
)

const (
	AuthTokenSize = 32
)

type Authentication struct {
	Id        string    `db:"id" json:"id"`
	UserId    string    `db:"user_id" json:"userId,omitempty"`
	EMail     string    `db:"email" json:"email,omitempty"`
	Password  string    `db:"password" json:"password,omitempty"`
	Token     string    `db:"token" json:"token,omitempty"`
	Enabled   bool      `db:"enabled" json:"enabled,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

func DefaultAuthentication() *Authentication {
	t := time.Now()
	return &Authentication{
		generateID(),
		"",
		"",
		"",
		generateRandomString(AuthTokenSize),
		false,
		t,
		t,
	}
}
