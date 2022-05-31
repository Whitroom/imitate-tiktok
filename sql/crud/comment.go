package crud

import (
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"gorm.io/gorm"
)

func CreateComment(db *gorm.DB, comment *models.Comment) *models.Comment {
	db.Create(&comment).Commit()
	return comment
}

func DeleteComment(db *gorm.DB, commentID uint) *models.Comment {
	var comment *models.Comment
	db.Model(&comment).Delete("id = ?", commentID).Commit()
	return comment
}

func GetComments(db *gorm.DB, videoID uint) []models.Comment {
	var comments []models.Comment
	db.Where("video_id = ?", videoID).Order("created_at desc").Find(&comments)
	return comments
}
