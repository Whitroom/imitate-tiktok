package controller

import (
	"net/http"

	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
)

type FavoriteActionRequest struct {
	VideoID    uint `form:"video_id" binding:"required"`
	ActionType uint `form:"action_type" binding:"required,min=1,max=2"`
}

func FavoriteAction(ctx *gin.Context) {
	var request FavoriteActionRequest
	if !BindAndValid(ctx, &request) {
		return
	}

	user_, _ := ctx.Get("User")
	user, _ := user_.(*models.User)

	if request.ActionType == 1 {
		if err := crud.UserLikeVideo(user.ID, request.VideoID); err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				StatusCode: 2,
				StatusMsg:  err.Error(),
			})
			return
		}
	} else {
		if err := crud.UserDislikeVideo(user.ID, request.VideoID); err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				StatusCode: 3,
				StatusMsg:  err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "操作成功",
	})

}

func FavoriteList(ctx *gin.Context) {

	user_, _ := ctx.Get("User")
	user, _ := user_.(*models.User)

	videos := crud.GetUserLikeVideosByUserID(user.ID)

	ctx.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: VideosModelChange(videos),
	})
}
