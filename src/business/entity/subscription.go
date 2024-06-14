package entity

import (
	"database/sql"
	"time"
)

const (
	UnlimitedPlan = "unlimited"
	VerifiedPlan  = "verified"
)

type Subscription struct {
	ID        int64        `db:"id" json:"id"`
	UserId    int64        `db:"user_id" json:"user_id"`
	Plan      string       `db:"plan" json:"plan"`
	StartDate time.Time    `db:"start_date" json:"start_date"`
	EndDate   time.Time    `db:"end_date" json:"end_date"`
	CreatedAt sql.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}

type SubscriptionParam struct {
	Plan string `json:"plan" validate:"oneof=unlimited verified"`
}
