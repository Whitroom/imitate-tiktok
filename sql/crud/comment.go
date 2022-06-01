package crud

import (
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
)

func CreateComment(comment *models.Comment) *models.Comment {
	sql.DB.Create(&comment).Commit()
	return comment
}

func DeleteComment(commentID uint) *models.Comment {
	var comment *models.Comment
	sql.DB.Model(&comment).Delete("id = ?", commentID).Commit()
	return comment
}

func GetComments(videoID uint) []models.Comment {
	var comments []models.Comment
	sql.DB.Where("video_id = ?", videoID).Order("created_at desc").Find(&comments)
	return comments
}
