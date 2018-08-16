package model

import (
	"time"
)

var (
	ItemStatusNotSold = "notsold"
	ItemStatusSold    = "sold"
	ItemStatusSoldOut = "soldout"
)

type Item struct {
	ID          string     `db:"id" json:"id"`
	UserID      *string    `db:"user_id" json:"userId,omitempty"`
	Label       *string    `db:"label" json:"label,omitempty"`
	Description *string    `db:"description" json:"description,omitempty"`
	Price       *int       `db:"price" json:"price,omitempty"`
	Status      *string    `db:"status" json:"status,omitempty"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

func DefaultItem() *Item {
	return &Item{
		generateID(),
		nil,
		nil,
		nil,
		nil,
		&ItemStatusNotSold,
		nil,
		nil,
	}
}
