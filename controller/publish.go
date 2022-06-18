package controller

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"
	"unicode/utf8"

	"gitee.com/Whitroom/imitate-tiktok/common"
	"gitee.com/Whitroom/imitate-tiktok/common/response"
	"gitee.com/Whitroom/imitate-tiktok/middlewares"
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
)

func Publish(ctx *gin.Context) {
	db := sql.GetSession()

	token := ctx.PostForm("token")
	userID, err := middlewares.Parse(ctx, token)
	if err != nil {
		return
	}

	data, err := ctx.FormFile("data")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: response.BADREQUEST,
			StatusMsg:  "文件获取错误: " + err.Error(),
		})
		return
	}

	if data.Filename[len(data.Filename)-3:] != "mp4" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: response.BADREQUEST,
			StatusMsg:  "不支持的文件格式",
		})
		return
	}

	title := ctx.PostForm("title")
	if title == "" || utf8.RuneCountInString(title) > 20 {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: response.BADREQUEST,
			StatusMsg:  "标题获取错误",
		})
		return
	}
	filename := filepath.Base(data.Filename)

	rand.Seed(time.Now().Unix())
	finalName := fmt.Sprintf("%d_%s", rand.Intn(100000000), filename)

	saveFile := filepath.Join("./public/", finalName)
	if err := ctx.SaveUploadedFile(data, saveFile); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			StatusCode: response.INTERNALERROR,
			StatusMsg:  err.Error(),
		})
		return
	}

	cmd := exec.Command("ffmpeg", "-i", "public/"+finalName,
		"-frames:v", "1", "-f", "image2",
		"public/covers/"+finalName[:len(finalName)-4]+".jpg")
	if err := cmd.Run(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		ctx.JSON(http.StatusInternalServerError, response.Response{
			StatusCode: response.INTERNALERROR,
			StatusMsg:  err.Error(),
		})
		return
	}

	crud.CreateVideo(db, &models.Video{
		AuthorID: userID,
		Title:    finalName,
	})

	ctx.JSON(http.StatusOK, response.Response{
		StatusCode: response.SUCCESS,
		StatusMsg:  finalName + " 上传成功",
	})
}

func PublishList(ctx *gin.Context) {
	db := sql.GetSession()

	user := common.GetUserFromCtx(ctx)
	videos := crud.GetUserPublishVideosByID(db, user.ID)
	responseVideos := common.VideosModelChange(db, user.ID, videos)
	ctx.JSON(http.StatusOK, response.VideoListResponse{
		Response: response.Response{
			StatusCode: response.SUCCESS,
			StatusMsg:  "获取成功",
		},
		VideoList: responseVideos,
	})
}
