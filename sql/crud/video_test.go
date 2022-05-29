package crud

import (
	"fmt"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	_ "gorm.io/driver/mysql"
	_ "gorm.io/gorm"
	"testing"
)

func TestCreateVideo(t *testing.T) {
	v := CreateVideo(db, &models.Video{AuthorID: 2, Title: "hello"})
	fmt.Println(v.Title)
}

func TestGetVideoUserFavoritesByID(t *testing.T) {
	f := GetVideoUserFavoritesByID(db, 2)
	for _, user := range f {
		fmt.Println(user.ID)
	}
}

func TestGetVideoCommentsByID(t *testing.T) {
	c := GetVideoCommentsByID(db, 2)
	for _, comment := range c {
		fmt.Println(comment.ID)
	}
}
