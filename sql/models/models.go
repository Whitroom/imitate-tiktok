package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string    `gorm:"size:10"`
	Password       string    `gorm:"size:40"`
	Content        string    `gorm:"size:50"`
	Videos         []Video   `gorm:"ForeignKey:AuthorID"`
	Comments       []Comment `gorm:"many2many:comments;joinForeignKey:UserID"`
	FavoriteVideos []Video   `gorm:"many2many:user_favorite_videos"`
	Subscribers    []User    `gorm:"joinForeignKey:SubscriberID;many2many:subscribes"`
	Followers      []User    `gorm:"joinForeignKey:UserID;many2many:subscribes"`
}

type Video struct {
	gorm.Model
	AuthorID      uint
	Title         string    `gorm:"size:30"`
	Author        User      `gorm:"reference:ID"`
	UserFavorites []User    `gorm:"many2many:user_favorite_videos"`
	Comments      []Comment `gorm:"many2many:Comment;joinForeignKey:VideoID"`
}

type Comment struct {
	gorm.Model
	UserID  uint
	VideoID uint
	Content string `gorm:"size:100"`
}
