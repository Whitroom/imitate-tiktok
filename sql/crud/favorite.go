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

	db.Model(&user).Association("FavoriteVideos").Append(&video)
	db.Commit()

	return nil
}

func UserDislikeVideo(db *gorm.DB, UserID uint, VideoID int) (*models.User, error) {
	var user *models.User
	var video *models.Video

	db.First(&user, UserID)
	db.First(&video, VideoID)

	if user == nil || video == nil {
		return nil, fmt.Errorf("找不到用户或视频")
	}

	db.Model(&user).Association("FavoriteVideos").Delete(&video)
	db.Commit()

	return user, nil
}

func GetUserLikeVideosByUserID(db *gorm.DB, UserID uint) []models.Video {
	var user *models.User
	db.Preload("FavoriteVideos").Find(&user, UserID)
	return user.FavoriteVideos
}

func GetVideoLikesCount(db *gorm.DB, VideoID uint) int64 {
	var video *models.Video
	var count int64
	db.Preload("UserFavorites").Find(&video, VideoID).Count(&count)
	return count
}

func IsUserFavoriteVideo(db *gorm.DB, userID, videoID uint) bool {
	var video *models.Video
	db.Raw("select * from video where id in "+
		"(select * from user_favorite_video where user_id = ? and video_id = ?)", userID, videoID)
	return video != nil
}
