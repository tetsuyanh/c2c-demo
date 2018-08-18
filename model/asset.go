package model

import (
	"time"
)

type Asset struct {
	Id        string     `db:"id" json:"id"`
	UserId    *string    `db:"user_id" json:"userId,omitempty"`
	Point     *int       `db:"point" json:"point,omitempty"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

func DefaultAsset() *Asset {
	p := 0
	t := time.Now()
	return &Asset{
		generateID(),
		nil,
		&p,
		&t,
		&t,
	}
}