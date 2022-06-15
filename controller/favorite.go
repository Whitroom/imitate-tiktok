package controller

import (
	"net/http"

	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"github.com/gin-gonic/gin"
)

func FavoriteAction(ctx *gin.Context) {
	var request struct {
		VideoID    uint `form:"video_id" binding:"required"`
		ActionType uint `form:"action_type" binding:"required,min=1,max=2"`
	}

	if !BindAndValid(ctx, &request) {
		return
	}

	user := GetUserFromCtx(ctx)

	if request.ActionType == 1 {
		if crud.IsUserFavoriteVideo(user.ID, request.VideoID) {
			ctx.JSON(http.StatusBadRequest, Response{
				StatusCode: 3,
				StatusMsg:  "已点赞过视频",
			})
			return
		}
		if err := crud.UserLikeVideo(user.ID, request.VideoID); err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				StatusCode: 2,
				StatusMsg:  err.Error(),
			})
			return
		}
	} else {

		if !crud.IsUserFavoriteVideo(user.ID, request.VideoID) {
			ctx.JSON(http.StatusBadRequest, Response{
				StatusCode: 3,
				StatusMsg:  "未点赞过视频",
			})
			return
		}
		if err := crud.UserDislikeVideo(user.ID, request.VideoID); err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				StatusCode: 2,
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

	userID := QueryIDAndValid(ctx, "user_id")

	videos := crud.GetUserLikeVideosByUserID(userID)

	ctx.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: VideosModelChange(videos),
	})
}
