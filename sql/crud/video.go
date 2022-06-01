package crud

import (
	"fmt"
	"time"

	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
)

func CreateVideo(video *models.Video) *models.Video {
	sql.DB.Create(&video).Commit()
	return video
}

func GetVideos(latestTime int64, userID uint) []models.Video {
	var videos []models.Video
	statement := sql.DB.Preload("Author").Limit(30)
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

func GetUserPublishVideosByID(userID uint) []models.Video {
	var videos []models.Video
	sql.DB.Preload("Author").Where("author_id = ?", userID).Find(&videos)
	return videos
}

func GetVideoByID(videoID uint) (*models.Video, error) {
	var video *models.Video
	sql.DB.First(&video, videoID)
	if video == nil {
		return nil, fmt.Errorf("未找到视频")
	}
	return video, nil
}

func GetVideoUserFavoritesByID(videoID uint) []models.User {
	var UserFavorites []models.User
	sql.DB.Raw("select * from users where id in("+
		"    select author_id from "+
		"videos where author_id= ?)", videoID).Scan(&UserFavorites)
	return UserFavorites
}

func GetVideoCommentsByID(videoID uint) []models.Comment {
	var Comments []models.Comment
	sql.DB.Raw("select * from comments"+
		" where video_id = ?", videoID).Scan(&Comments)
	return Comments
}

func GetVideoCommentsCountByID(videoID uint) int64 {
	var count int64
	sql.DB.Raw("select count(*) from comments"+
		" where video_id = ?", videoID).Scan(&count)
	return count
}
