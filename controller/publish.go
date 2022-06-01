package controller

import (
	"fmt"
	"net/http"
	"path/filepath"

	"gitee.com/Whitroom/imitate-tiktok/middlewares"
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
	userID, err := middlewares.Parse(ctx, token)
	if err != nil {
		return
	}

	data, err := ctx.FormFile("data")
	if err != nil {
		ctx.JSON(http.StatusOK, Response{
			StatusCode: 2,
			StatusMsg:  "文件获取错误: " + err.Error(),
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
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 2,
			StatusMsg:  "标题获取错误",
		})
	}
	filename := filepath.Base(data.Filename)

	video := crud.CreateVideo(&models.Video{
		AuthorID: userID,
		Title:    title,
	})

	finalName := fmt.Sprintf("%d_%s", video.ID, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := ctx.SaveUploadedFile(data, saveFile); err != nil {
		ctx.JSON(http.StatusOK, Response{
			StatusCode: 4,
			StatusMsg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " 上传成功",
	})
}

func PublishList(ctx *gin.Context) {
	user_, _ := ctx.Get("User")
	user, _ := user_.(*models.User)
	videos := crud.GetUserPublishVideosByID(user.ID)
	modelVideos := VideosModelChange(videos)
	for i := 0; i < len(modelVideos); i++ {
		modelVideos[i].IsFavorite = crud.IsUserFavoriteVideo(user.ID, uint(modelVideos[i].Id))
	}
	ctx.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: modelVideos,
	})
}
