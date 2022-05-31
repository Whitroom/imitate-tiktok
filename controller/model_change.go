package controller

import (
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
)

func UsersModelChange(Users []models.User) []User {
	var users []User
	for _, user := range Users {
		users = append(users, UserModelChange(user))
	}
	return users
}

func UserModelChange(user models.User) User {
	return User{
		Id:            int64(user.ID),
		Name:          user.Name,
		FollowCount:   crud.GetUserSubscribersCountByID(sql.DB, user.ID),
		FollowerCount: crud.GetUserFollowersCountByID(sql.DB, user.ID),
		IsFollow:      true,
	}
}

func CommentsModelChange(Comments []models.Comment) []Comment {
	var comments []Comment
	for _, comment := range Comments {
		comments = append(comments, CommentModelChange(comment))
	}
	return comments
}

func CommentModelChange(comment models.Comment) Comment {
	user, _ := crud.GetUserByID(sql.DB, comment.UserID)
	return Comment{
		Id:         int64(comment.ID),
		Content:    comment.Content,
		CreateDate: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		User:       UserModelChange(*user),
	}
}

func VideosModelChange(Videos []models.Video) []Video {
	var videos []Video
	for _, video := range Videos {
		videos = append(videos, VideoModelChange(&video))
	}
	return videos
}

func VideoModelChange(video *models.Video) Video {
	return Video{
		Id:            int64(video.ID),
		FavoriteCount: crud.GetVideoLikesCount(sql.DB, video.ID),
		Author:        UserModelChange(video.Author),
		CommentCount:  crud.GetVideoCommentsCountByID(sql.DB, video.ID),
		IsFavorite:    true,
		// 以下是测试数据
		PlayUrl:  "http://your_machine_ip:8080/static/" + video.Title,
		CoverUrl: "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
	}
}
