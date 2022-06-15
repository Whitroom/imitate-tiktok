package crud

import (
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"gorm.io/gorm"
)

func CreateComment(db *gorm.DB, comment *models.Comment) *models.Comment {
	sql.DB.Create(&comment).Commit()
	return comment
}

func DeleteComment(db *gorm.DB, commentID uint) {
	var comment *models.Comment
	sql.DB.Model(&comment).Delete("id = ?", commentID).Commit()
}

func GetComments(db *gorm.DB, videoID uint) []models.Comment {
	var comments []models.Comment
	sql.DB.Where("video_id = ?", videoID).Order("created_at desc").Find(&comments)
	return comments
}

func GetVideoCommentsCountByID(db *gorm.DB, videoID uint) int64 {
	var count int64
	sql.DB.Raw("select count(*) from comments"+
		" where video_id = ? and deleted_at is not null", videoID).Scan(&count)
	return count
}
