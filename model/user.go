package model

import "database/sql"

type User struct {
	Id       int    `json:"-" `
	Email    string `json:"email" default:"babdusalom72@gmail.com"`
	Password string `json:"password" default:"1996"`
	UserName string `json:"username" default:"MY"`
	Image    string `json:"image" `
	City     string `json:"city" default:"Navoi"`
	Phone    string `json:"phone" default:"+9989 99 331 96 12"`
}
type UserUpdate struct {
	UserName string `json:"username" default:"MY"`

	City  string `json:"city" default:"Navoi"`
	Phone string `json:"phone" default:"+9989 99 331 96 12"`
}

type UserFull struct {
	Id                int            `json:"-" db:"id"`
	Email             string         `json:"email"  db:"email"`
	UserName          string         `json:"username" db:"username"`
	City              string         `json:"city" db:"city"`
	IsVerified        bool           `json:"is_verified" db:"is_verified"`
	Verification_date sql.NullTime   `json:"verification_date" db:"verification_date"`
	AccountImagePath  sql.NullString `json:"account_image_path" db:"account_image_path"`
	Phone             string         `json:"phone" db:"phone"`
	Rating            string         `json:"rating" db:"rating"`
	PostViews         int            `json:"post_views" db:"post_views"`
	IsSuperUser       bool           `json:"is_super_user" db:"is_super_user"`
	CreatedAt         string         `json:"created_at" db:"created_at"`
}
type ResponseSignUp struct {
	Id    int    `json:"id" `
	Token string `json:"token"`
}
type ResponseSignIn struct {
	Email string `json:"email" `
	Token string `json:"token"`
}

type SignInInput struct {
	Email    string `json:"email"  default:"babdusalom72@gmail.com"`
	Password string `json:"password" default:"1996"`
}

type VerificationCode struct {
	Code string `json:"code"`
}
