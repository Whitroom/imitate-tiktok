package crud

import (
	"fmt"

	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"gorm.io/gorm"
)

func CreateVideo(db *gorm.DB, video *models.Video) *models.Video {
	db.Create(&video).Commit()
	return video
}

func GetVideoByID(db *gorm.DB, videoID uint) (*models.Video, error) {
	var video *models.Video
	db.First(&video, videoID)
	if video == nil {
		return nil, fmt.Errorf("未找到视频")
	}
	return video, nil
}

func GetVideoUserFavoritesByID(db *gorm.DB, videoID uint) []models.User {
	var UserFavorites []models.User
	db.Raw("select * from users where id in("+
		"    select author_id from "+
		"videos where author_id= ?)", videoID).Scan(&UserFavorites)
	return UserFavorites
}

func GetVideoCommentsByID(db *gorm.DB, videoID uint) []models.Comment {
	var Comments []models.Comment
	db.Raw("select * from comments"+
		" where video_id = ?", videoID).Scan(&Comments)
	return Comments
}

func GetVideoCommentsCountByID(db *gorm.DB, videoID uint) int64 {
	var count int64
	db.Raw("select * from comments"+
		" where video_id = ?", videoID).Count(&count)
	return count
}
