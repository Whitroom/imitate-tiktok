package controller

import (
	"net/http"

	"gitee.com/Whitroom/imitate-tiktok/common"
	"gitee.com/Whitroom/imitate-tiktok/common/response"
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"github.com/gin-gonic/gin"
)

func FavoriteAction(ctx *gin.Context) {
	db := sql.GetDB()

	var request struct {
		VideoID    uint `form:"video_id" binding:"required"`
		ActionType uint `form:"action_type" binding:"required,min=1,max=2"`
	}

	if !common.BindAndValid(ctx, &request) {
		return
	}

	user := common.GetUserFromCtx(ctx)

	if request.ActionType == 1 {
		if crud.IsUserFavoriteVideo(&db, user.ID, request.VideoID) {
			ctx.JSON(http.StatusBadRequest, response.Response{
				StatusCode: 3,
				StatusMsg:  "已点赞过视频",
			})
			return
		}
		if err := crud.UserLikeVideo(&db, user.ID, request.VideoID); err != nil {
			ctx.JSON(http.StatusNotFound, response.Response{
				StatusCode: 2,
				StatusMsg:  err.Error(),
			})
			return
		}
	} else {

		if !crud.IsUserFavoriteVideo(&db, user.ID, request.VideoID) {
			ctx.JSON(http.StatusBadRequest, response.Response{
				StatusCode: 3,
				StatusMsg:  "未点赞过视频",
			})
			return
		}
		if err := crud.UserDislikeVideo(&db, user.ID, request.VideoID); err != nil {
			ctx.JSON(http.StatusNotFound, response.Response{
				StatusCode: 2,
				StatusMsg:  err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, response.Response{
		StatusCode: 0,
		StatusMsg:  "操作成功",
	})

}

func FavoriteList(ctx *gin.Context) {
	db := sql.GetDB()

	userID := common.QueryIDAndValid(ctx, "user_id")

	videos := crud.GetUserLikeVideosByUserID(&db, userID)

	ctx.JSON(http.StatusOK, response.VideoListResponse{
		Response: response.Response{
			StatusCode: 0,
		},
		VideoList: common.VideosModelChange(&db, videos),
	})
}
