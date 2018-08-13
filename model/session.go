package model

import "time"

// SessionTokenSize is byte size of SessionTokenSize
const SessionTokenSize = 32

// Session has rights to behave as a paticular user
type Session struct {
	ID        string     `db:"id" json:"id"`
	UserID    string     `db:"user_id" json:"userId,omitempty"`
	Token     string     `db:"token" json:"token"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt,omitempty"`
}

func NewSession(userID string) *Session {
	t := time.Now()
	return &Session{
		generateID(),
		userID,
		generateRandomString(SessionTokenSize),
		&t,
	}
}
