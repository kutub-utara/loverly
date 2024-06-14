package entity

import (
	"database/sql"
	"time"
)

const (
	Male   = "male"
	Female = "female"
)

type Profile struct {
	ID        int64          `db:"id" json:"id"`
	UserId    int64          `db:"user_id" json:"user_id"`
	FullName  string         `db:"name" json:"fullname"`
	BirthDay  sql.NullTime   `db:"birthday" json:"-"`
	Age       int64          `db:"-" json:"age"`
	Gender    string         `db:"gender" json:"gender"`
	Location  sql.NullString `db:"location" json:"location"`
	Bio       sql.NullString `db:"bio" json:"bio"`
	ProfPic   sql.NullString `db:"profile_picture" json:"profile_picture"`
	Interest  sql.NullString `db:"interests" json:"interests"`
	CreatedAt sql.NullTime   `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at" json:"deleted_at"`
}

type ProfileResponse struct {
	FullName  string    `json:"fullname"`
	Age       int64     `json:"age"`
	Gender    string    `json:"gender"`
	Location  string    `json:"location"`
	Bio       string    `json:"bio"`
	ProfPic   string    `json:"profile_picture"`
	Interest  string    `json:"interests"`
	CreatedAt time.Time `json:"created_at"`
}
