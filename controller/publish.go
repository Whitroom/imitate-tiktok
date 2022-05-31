package controller

import (
	"fmt"
	"net/http"
	"path/filepath"

	"gitee.com/Whitroom/imitate-tiktok/middlewares"
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

func Publish(ctx *gin.Context) {

	token := ctx.PostForm("token")
	userID, err := middlewares.Parse(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    2,
			"message": "token获取错误, 请重新登陆获取",
		})
		return
	}

	data, err := ctx.FormFile("data")
	if err != nil {
		ctx.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	if data.Filename[len(data.Filename)-3:] != "mp4" {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 3,
			StatusMsg:  "不支持的文件格式",
		})
		return
	}

	title := ctx.PostForm("title")
	fmt.Println(title)
	if title == "" {
		ctx.JSON(http.StatusBadRequest, Response{StatusCode: 2, StatusMsg: "Title not Found"})
	}
	filename := filepath.Base(data.Filename)

	finalName := fmt.Sprintf("%d_%s_%s", userID, title, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := ctx.SaveUploadedFile(data, saveFile); err != nil {
		ctx.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	crud.CreateVideo(sql.DB, &models.Video{
		AuthorID: userID,
		Title:    finalName,
	})

	ctx.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

func PublishList(ctx *gin.Context) {
	user_, _ := ctx.Get("User")
	user, _ := user_.(*models.User)
	videos := crud.GetUserPublishVideosByID(sql.DB, user.ID)
	modelVideos := VideosModelChange(videos)
	for i := 0; i < len(modelVideos); i++ {
		modelVideos[i].IsFavorite = crud.IsUserFavoriteVideo(sql.DB, user.ID, uint(modelVideos[i].Id))
	}
	ctx.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: modelVideos,
	})
}
