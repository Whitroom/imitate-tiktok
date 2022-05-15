package sql

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	ShowLikes   bool     `json:"show_likes"`
	Enabled     bool     `json:"enabled"`
	Subscribers []*User  `gorm:"many2many:subscribes"`
	Followers   []*User  `gorm:"many2many:subscribes"`
	Likes       []*Video `gorm:"many2many:likes"`
	Comments    []*Video `gorm:"many2many:comments"`
}

type Subscribe struct {
	gorm.Model
	SubscriberID int `json:"producer_id"`
	FollowerID   int `json:"follower_id"`
}

type Video struct {
	gorm.Model
	VideoStatus string `json:"video_status"`
	Description string `json:"description"`
}

type Like struct {
	gorm.Model
	UserID  int  `json:"user_id"`
	VideoID int  `json:"video_id"`
	Enabled bool `json:"enabled"`
}

type Comment struct {
	gorm.Model
	UserID  int    `json:"user_id"`
	VideoID int    `json:"video_id"`
	Context string `json:"context"`
	Enabled bool   `json:"enabled"`
}
