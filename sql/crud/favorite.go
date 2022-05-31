package crud

import (
	"fmt"

	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"gorm.io/gorm"
)

func UserLikeVideo(db *gorm.DB, UserID uint, VideoID uint) error {
	var user *models.User
	var video *models.Video

	db.First(&user, UserID)
	db.First(&video, VideoID)

	if user == nil || video == nil {
		return fmt.Errorf("找不到用户或视频")
	}

	db.Model(&user).Association("FavoriteVideos").Append(video)
	db.Commit()

	return nil
}

func UserDislikeVideo(db *gorm.DB, UserID uint, VideoID uint) error {
	var user *models.User
	var video *models.Video

	db.First(&user, UserID)
	db.First(&video, VideoID)

	if user == nil || video == nil {
		return fmt.Errorf("找不到用户或视频")
	}

	err := db.Model(&user).Association("FavoriteVideos").Delete(video)
	if err != nil {
		return fmt.Errorf("找不到点赞的视频")
	}
	db.Commit()

	return nil
}

func GetUserLikeVideosByUserID(db *gorm.DB, UserID uint) []models.Video {
	var user *models.User
	db.Preload("FavoriteVideos").Find(&user, UserID)
	return user.FavoriteVideos
}

func GetVideoLikesCount(db *gorm.DB, VideoID uint) int64 {
	var count int64
	db.Raw("select count(user_id) from user_favorite_videos where video_id = ?", VideoID).Scan(&count)
	return count
}

func IsUserFavoriteVideo(db *gorm.DB, userID, videoID uint) bool {
	var video *models.Video
	if userID == 0 {
		return false
	}
	db.Raw("select * from videos where id in "+
		"(select video_id from user_favorite_videos where user_id = ? and video_id = ?)",
		userID, videoID).Scan(&video)
	return video != nil
}
