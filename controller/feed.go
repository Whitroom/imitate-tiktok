package controller

import (
	"net/http"
	"strconv"

	"gitee.com/Whitroom/imitate-tiktok/middlewares"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// 如果出现token 则不会出现自己的视频
func Feed(ctx *gin.Context) {
	var latestTime, nextTime int64
	token := ctx.Query("token")
	latestTime_ := ctx.Query("latest_time")
	if latestTime_ != "" {
		latestTime, _ = strconv.ParseInt(latestTime_, 10, 64)
	} else {
		latestTime = 0
	}
	var userID uint
	if token != "" {
		var err error
		userID, err = middlewares.Parse(ctx, token)
		if err != nil {
			return
		}

	} else {
		userID = 0
	}
	videos := crud.GetVideos(latestTime, userID)
	modelVideos := VideosModelChange(videos)
	for i := 0; i < len(modelVideos); i++ {
		modelVideos[i].IsFavorite = crud.IsUserFavoriteVideo(userID, uint(modelVideos[i].Id))
	}
	if len(videos)-1 < 0 {
		nextTime = 0
	} else {
		nextTime = videos[len(videos)-1].CreatedAt.Unix()
	}
	ctx.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: modelVideos,
		NextTime:  nextTime,
	})
}
