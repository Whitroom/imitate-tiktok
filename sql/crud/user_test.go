package crud

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var (
	db = GetDB()
)

func TestGetVideoByAuthorID(t *testing.T) {
	g := GetVideoByAuthorID(db, 1)
	for _, video := range g {
		fmt.Println(video.ID)
	}
}

func TestGetUserFollowersById(t *testing.T) {
	got := GetUserFollowersById(db, 1)
	for _, user := range got {
		fmt.Println(user.ID)
	}
}
func GetDB() *gorm.DB {
	dsn := "root:111111@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db
}
