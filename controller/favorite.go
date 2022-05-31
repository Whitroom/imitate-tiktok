package controller

import (
	"net/http"

	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
)

type FavoriteActionRequest struct {
	// UserID     uint `form:"user_id" binding:"required"`
	VideoID    uint `form:"video_id" binding:"required"`
	ActionType uint `form:"action_type" binding:"required,min=1,max=2"`
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(ctx *gin.Context) {
	var request FavoriteActionRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "数据绑定失败",
		})
		return
	}
	user_, _ := ctx.Get("User")
	user, _ := user_.(*models.User)

	if request.ActionType == 1 {
		if err := crud.UserLikeVideo(sql.DB, user.ID, request.VideoID); err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				StatusCode: 2,
				StatusMsg:  err.Error(),
			})
			return
		}
	} else {
		if err := crud.UserDislikeVideo(sql.DB, user.ID, request.VideoID); err != nil {
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

// FavoriteList all users have same favorite video list
func FavoriteList(ctx *gin.Context) {

	user, _ := ctx.Get("User")
	user_, _ := user.(*models.User)

	videos := crud.GetUserLikeVideosByUserID(sql.DB, user_.ID)

	ctx.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: VideosModelChange(videos),
	})
}
