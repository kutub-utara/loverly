package entity

import (
	"database/sql"

	"golang.org/x/oauth2"
)

const (
	NextStateLogin  = "login"
	NextStateVerify = "verify"
)

type User struct {
	ID        int64        `db:"id"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	Verifed   bool         `db:"verified"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type SignInParam struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignInResponse struct {
	ID      int64         `json:"id"`
	Email   string        `json:"email"`
	Verifed bool          `json:"verified"`
	Token   *oauth2.Token `json:"token"`
}

type SignUpParam struct {
	FullName        string `json:"fullname" validate:"required"`
	Gender          string `json:"gender" validate:"oneof=male female"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}

type SignUpResponse struct {
	NextState string `json:"next_state"`
}
