package common

import (
	"gitee.com/Whitroom/imitate-tiktok/common/response"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"gorm.io/gorm"
)

func UsersModelChange(db *gorm.DB, userList []models.User) []response.User {
	var users []response.User
	for _, user := range userList {
		users = append(users, UserModelChange(db, user))
	}
	return users
}

func UserModelChange(db *gorm.DB, user models.User) response.User {
	return response.User{
		ID:            int64(user.ID),
		Name:          user.Name,
		FollowCount:   crud.GetUserSubscribersCountByID(db, user.ID),
		FollowerCount: crud.GetUserFollowersCountByID(db, user.ID),
		IsFollow:      true,
	}
}

func CommentsModelChange(db *gorm.DB, commentList []models.Comment) []response.Comment {
	var comments []response.Comment
	for _, comment := range commentList {
		comments = append(comments, CommentModelChange(db, comment))
	}
	return comments
}

func CommentModelChange(db *gorm.DB, comment models.Comment) response.Comment {
	user, _ := crud.GetUserByID(db, comment.UserID)
	return response.Comment{
		ID:         int64(comment.ID),
		Content:    comment.Content,
		CreateDate: comment.CreatedAt.Format("01-02"),
		User:       UserModelChange(db, *user),
	}
}

func VideosModelChange(db *gorm.DB, userID uint, videoList []models.Video) []response.Video {
	var videos []response.Video
	for _, video := range videoList {
		videos = append(videos, VideoModelChange(db, userID, &video))
	}
	return videos
}

func VideoModelChange(db *gorm.DB, userID uint, video *models.Video) response.Video {
	return response.Video{
		ID:            int64(video.ID),
		FavoriteCount: crud.GetVideoLikesCount(db, video.ID),
		Author:        UserModelChange(db, video.Author),
		CommentCount:  crud.GetVideoCommentsCountByID(db, video.ID),
		IsFavorite:    crud.IsUserFavoriteVideo(db, userID, video.ID),
		// 以下是测试数据
		PlayUrl:  "http://192.168.1.4:8080/static/" + video.Title,
		CoverUrl: "http://192.168.1.4:8080/static/covers/" + video.Title[:len(video.Title)-3] + "jpg",
	}
}
