package controller

import (
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
)

func UserModelChange(user models.User) User {
	return User{
		Id:            int64(user.ID),
		Name:          user.Name,
		FollowCount:   crud.GetUserSubscribersCountByID(sql.DB, user.ID),
		FollowerCount: crud.GetUserFollowersCountByID(sql.DB, user.ID),
		IsFollow:      true,
	}
}

func UserPointerModelChange(user *models.User) User {
	return User{
		Id:            int64(user.ID),
		Name:          user.Name,
		FollowCount:   crud.GetUserSubscribersCountByID(sql.DB, user.ID),
		FollowerCount: crud.GetUserFollowersCountByID(sql.DB, user.ID),
		IsFollow:      true,
	}
}

func CommentModelChange(comment models.Comment) Comment {
	user, _ := crud.GetUserByID(sql.DB, comment.UserID)
	return Comment{
		Id:         int64(comment.ID),
		Content:    comment.Content,
		CreateDate: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		User:       UserPointerModelChange(user),
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
		PlayUrl:  "",
		CoverUrl: "",
	}
}
