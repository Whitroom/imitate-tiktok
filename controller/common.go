package controller

import (
	"net/http"
	"strconv"

	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	ID            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,"`
	CommentCount  int64  `json:"comment_count,"`
	IsFavorite    bool   `json:"is_favorite"`
}

type Comment struct {
	ID         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_dat,"`
}

type User struct {
	ID            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func BindAndValid(ctx *gin.Context, target interface{}) bool {
	if err := ctx.ShouldBindQuery(target); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, Response{
			StatusCode: 1,
			StatusMsg:  "参数匹配错误",
		})
		return false
	}
	return true
}

func QueryIDAndValid(ctx *gin.Context, queryName string) uint {
	id, err := strconv.ParseUint(ctx.Query(queryName), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  queryName + "不是数字",
		})
		return 0
	}
	return uint(id)
}

func GetUserFromCtx(ctx *gin.Context) *models.User {
	user_, _ := ctx.Get("User")
	user, _ := user_.(*models.User)
	return user
}

func UsersModelChange(db *gorm.DB, userList []models.User) []User {
	var users []User
	for _, user := range userList {
		users = append(users, UserModelChange(db, user))
	}
	return users
}

func UserModelChange(db *gorm.DB, user models.User) User {
	return User{
		ID:            int64(user.ID),
		Name:          user.Name,
		FollowCount:   crud.GetUserSubscribersCountByID(db, user.ID),
		FollowerCount: crud.GetUserFollowersCountByID(db, user.ID),
		IsFollow:      true,
	}
}

func CommentsModelChange(db *gorm.DB, commentList []models.Comment) []Comment {
	var comments []Comment
	for _, comment := range commentList {
		comments = append(comments, CommentModelChange(db, comment))
	}
	return comments
}

func CommentModelChange(db *gorm.DB, comment models.Comment) Comment {
	user, _ := crud.GetUserByID(db, comment.UserID)
	return Comment{
		ID:         int64(comment.ID),
		Content:    comment.Content,
		CreateDate: comment.CreatedAt.Format("01-02"),
		User:       UserModelChange(db, *user),
	}
}

func VideosModelChange(db *gorm.DB, videoList []models.Video) []Video {
	var videos []Video
	for _, video := range videoList {
		videos = append(videos, VideoModelChange(db, &video))
	}
	return videos
}

func VideoModelChange(db *gorm.DB, video *models.Video) Video {
	return Video{
		ID:            int64(video.ID),
		FavoriteCount: crud.GetVideoLikesCount(db, video.ID),
		Author:        UserModelChange(db, video.Author),
		CommentCount:  crud.GetVideoCommentsCountByID(db, video.ID),
		IsFavorite:    true,
		// 以下是测试数据
		PlayUrl:  "http://your-machine-ip:8080/static/" + video.Title,
		CoverUrl: "http://your-machine-ip:8080/static/covers/" + video.Title[:len(video.Title)-4] + ".jpg",
	}
}
