package crud

import (
	"fmt"
	"time"

	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"gorm.io/gorm"
)

func CreateVideo(db *gorm.DB, video *models.Video) *models.Video {
	db.Create(&video).Commit()
	return video
}

func GetVideos(db *gorm.DB, latestTime int64, userID uint) []models.Video {
	var videos []models.Video
	statement := db.Preload("Author").Limit(30)
	if latestTime != 0 {
		statement = statement.Where("created_at < ?",
			time.Unix(latestTime/1000+43200, 0).Format("2006-01-02 15:04:05"))
	}
	if userID != 0 {
		statement = statement.Where("author_id != ?", userID)
	}
	statement.Order("created_at desc").Find(&videos)
	return videos
}

func GetUserPublishVideosByID(db *gorm.DB, userID uint) []models.Video {
	var videos []models.Video
	db.Preload("Author").Where("author_id = ?", userID).Find(&videos)
	return videos
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
	db.Raw("select count(*) from comments"+
		" where video_id = ?", videoID).Scan(&count)
	return count
}
