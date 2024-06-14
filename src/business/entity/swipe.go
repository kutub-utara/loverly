package entity

import "database/sql"

type Swipe struct {
	ID        int64        `db:"id" json:"id"`
	SwiperId  int64        `db:"swiper_id" json:"swiper_id"`
	SwipedId  int64        `db:"swiped_id" json:"swiped_id"`
	Direction string       `db:"direction" json:"direction"`
	CreatedAt sql.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}

type SwipeParam struct {
	SwipedId  int64  `json:"swiped_id" validate:"required"`
	Direction string `json:"direction" validate:"oneof=left right"`
}

type SwipeResponse struct {
	Match bool `json:"match,omitempty"`
	Like  bool `json:"like,omitempty"`
}
