package entity

import "database/sql"

type Match struct {
	ID        int64        `db:"id" json:"id"`
	UserId1   int64        `db:"user_id_1" json:"user_id_1"`
	UserId2   int64        `db:"user_id_2" json:"user_id_2"`
	CreatedAt sql.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}
