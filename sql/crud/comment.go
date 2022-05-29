package crud

import (
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"gorm.io/gorm"
)

func CreateComment(db *gorm.DB, comment *models.Comment) *models.Comment {
	db.Create(&comment).Commit()
	return comment
}
