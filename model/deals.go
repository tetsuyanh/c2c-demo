package model

import (
	"time"
)

type Deal struct {
	Id        string    `db:"id" json:"id"`
	ItemId    string    `db:"item_id" json:"itemId,omitempty"`
	SellerId  string    `db:"seller_id" json:"sellerId,omitempty"`
	BuyerId   string    `db:"buyer_id" json:"buyerId,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"createdAt,omitempty"`
}

func DefaultDeal() *Deal {
	return &Deal{
		generateID(),
		"",
		"",
		"",
		time.Now(),
	}
}
