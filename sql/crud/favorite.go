package crud

import (
	"fmt"

	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"gorm.io/gorm"
)

func UserLikeVideo(db *gorm.DB, userID uint, videoID uint) error {
	var user *models.User
	var video *models.Video

	db.First(&user, userID)
	db.First(&video, videoID)

	if user == nil || video == nil {
		return fmt.Errorf("找不到用户或视频")
	}

	db.Model(&user).Association("FavoriteVideos").Append(video)
	db.Commit()

	return nil
}

func UserDislikeVideo(db *gorm.DB, userID uint, videoID uint) error {
	var user *models.User
	var video *models.Video

	db.First(&user, userID)
	db.First(&video, videoID)

	if user == nil || video == nil {
		return fmt.Errorf("找不到用户或视频")
	}

	if db.Model(&user).Association("FavoriteVideos").Delete(video) != nil {
		return fmt.Errorf("找不到点赞的视频")
	}
	db.Commit()

	return nil
}

func GetUserLikeVideosByUserID(db *gorm.DB, userID uint) []models.Video {
	var videos []models.Video
	db.Raw("select * from videos where id in (select video_id from user_favorite_videos where user_id = ?)", userID).Scan(&videos)
	return videos
}

func GetVideoLikesCount(db *gorm.DB, videoID uint) int64 {
	var count int64
	db.Raw("select count(user_id) from user_favorite_videos where video_id = ?", videoID).Scan(&count)
	return count
}

func IsUserFavoriteVideo(db *gorm.DB, userID, videoID uint) bool {
	if userID == 0 {
		return false
	}
	err := db.Raw("select * from videos where id in "+
		"(select video_id from user_favorite_videos where user_id = ? and video_id = ?)",
		userID, videoID).Scan(&models.Video{}).Error
	return err == nil
}
