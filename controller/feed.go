package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	//var videoList []Video

	//get videoName form public directory
	//FilesName, err := GetAllFile("public", []string{})
	//if err != nil {
	//	fmt.Printf("file error %w", err)
	//	return
	//}

	//
	//for _, fileName := range FilesName {
	//	video:=Video{
	//		Author: ,
	//		PlayUrl: fmt.Sprintf("http://%v:8080/static/%v",,fileName),
	//		FavoriteCount: ,
	//		CommentCount: ,
	//		IsFavorite: ,
	//	}
	//}

	//c.JSON(http.StatusOK, FeedResponse{
	//	Response:  Response{StatusCode: 0},
	//	VideoList: videoList,
	//	NextTime:  time.Now().Unix(),
	//})

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}
