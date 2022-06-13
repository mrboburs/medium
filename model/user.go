package model

import (
	"database/sql"
	// "github.com/lib/pq"
)

type User struct {
	Id        int    `json:"-" `
	Email     string `json:"email" default:"babdusalom72@gmail.com"`
	Password  string `json:"password" default:"1996"`
	FirstName string `json:"first_name" default:"MY"`
	LastName  string `json:"last_name" default:"mrb"`
	Interests string `json:"interests" `
	Bio       string `json:"bio" default:"I am Golang dev"`
	City      string `json:"city" default:"Navoi"`
	Phone     string `json:"phone" default:"+9989 99 331 96 12"`
}

type UserUpdate struct {
	FirstName string `json:"first_name" default:"MY"`
	LastName  string `json:"last_name" default:"mrb"`
	Interests string `json:"interests" `
	Bio       string `json:"bio" default:"I am Golang dev"`
	City      string `json:"city" default:"Navoi"`
	Phone     string `json:"phone" default:"+9989 99 331 96 12"`
}

type UserFull struct {
	Id                int          `json:"id" db:"id"`
	Email             string       `json:"email"  db:"email"`
	FirstName         string       `json:"first_name" db:"first_name"`
	Lastname          string       `json:"last_name" db:"last_name"`
	City              string       `json:"city" db:"city"`
	IsVerified        bool         `json:"is_verified" db:"is_verified"`
	Verification_date sql.NullTime `json:"verification_date" db:"verification_date"`
	Interests         string       `json:"interests" db:"interests"`
	Bio               string       `json:"bio" db:"bio"`
	AccountImagePath  string       `json:"account_image_path" db:"account_image_path"`
	Phone             string       `json:"phone" db:"phone"`
	Rating            string       `json:"rating" db:"rating"`
	PostViewsCount    int          `json:"post_views" db:"post_views_count"`
	FollowerCount     int          `json:"follower_count" db:"follower_count"`
	FollowingCount    int          `json:"following_count" db:"following_count"`
	LikeCount         int          `json:"like_count" db:"like_count"`
	IsSuperUser       bool         `json:"is_super_user" db:"is_super_user"`
	CreatedAt         string       `json:"created_at" db:"created_at"`
}
