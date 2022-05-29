package controller

import (
	"errors"
	"fmt"
	"gitee.com/Whitroom/imitate-tiktok/middlewares"
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.Query("token")

	//if _, exist := usersLoginInfo[token]; !exist {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//	return
	//}

	//get file and title
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "File not Found",
		})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, Response{StatusCode: 2, StatusMsg: "Title not Found"})
	}

	filename := filepath.Base(data.Filename)
	userID, err := middlewares.Parse(token)
	user, err := crud.GetUserByID(sql.DB, userID)

	//save video to database
	video := crud.CreateVideo(sql.DB, &models.Video{AuthorID: userID, Author: *user, Title: title})

	//Save video file to pulbic directory
	finalName := fmt.Sprintf("%d_%s", video.ID, filename)
	saveFile := filepath.Join("./public/", finalName)

	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userID := c.Query("user_id")
	if userID != "" {
		fmt.Println("userID is empty")
		return
	}
	//get video list by user_id
	//type Video struct {
	//	Id            int64  `json:"id,omitempty"`
	//	Author        User   `json:"author"`
	//	PlayUrl       string `json:"play_url,omitempty"`
	//	CoverUrl      string `json:"cover_url,omitempty"`
	//	FavoriteCount int64  `json:"favorite_count,omitempty"`
	//	CommentCount  int64  `json:"comment_count,omitempty"`
	//	IsFavorite    bool   `json:"is_favorite,omitempty"`
	//}

	//uID, _ := strconv.Atoi(userID)
	//videos := crud.GetVideoByAuthorID(sql.DB, uint(uID))
	//var videoList []Video
	//for _, modelsVideo := range videos {
	//	video:=Video{
	//		Id: int64(uint(modelsVideo.ID)),
	//		Author: modelsVideo.Author,
	//		PlayUrl: ,
	//
	//
	//	}
	//}
	//DemoVideos
	//	{
	//		Id:            2,
	//		Author:        DemoUser,
	//		PlayUrl:       "http://192.168.50.45:8080/static/2_ai.mp4",
	//		FavoriteCount: 0,
	//		CommentCount:  0,
	//		IsFavorite:    false,
	//	},

	c.JSON(http.StatusOK, VideoListResponse{

		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}

//GetAllFile get video file from pathname directory(supporting mp4,avi,wmv)
func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	videoFormat := []string{"mp4", "avi", "wmv"}
	for _, fi := range rd {
		if fi.IsDir() == false {
			ext := filepath.Ext(fi.Name())
			ext = ext[1:]
			err := ContainElement(ext, videoFormat)
			if err != nil {
				continue
			}
			s = append(s, ext)
		}
	}
	return s, nil
}

//ContainElement whether support this video format
func ContainElement(ext string, videoFormat []string) error {
	for _, vf := range videoFormat {
		if ext == vf {
			return nil
		}
	}
	return errors.New("not support this format")
}
