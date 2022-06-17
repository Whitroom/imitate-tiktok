package controller

import (
	"net/http"

	"gitee.com/Whitroom/imitate-tiktok/common"
	"gitee.com/Whitroom/imitate-tiktok/common/response"
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
)

func CommentAction(ctx *gin.Context) {
	db := sql.GetSession()

	var request struct {
		VideoID     uint   `form:"video_id" binding:"required"`
		ActionType  uint   `form:"action_type" binding:"required,min=1,max=2"`
		CommentText string `form:"comment_text" binding:"omitempty"`
		CommentID   uint   `form:"comment_id" binding:"omitempty"`
	}

	if !common.BindAndValid(ctx, &request) {
		return
	}
	user := common.GetUserFromCtx(ctx)
	if request.ActionType == 1 {
		if len(request.CommentText) == 0 {
			ctx.JSON(http.StatusBadRequest, response.Response{
				StatusCode: response.BADREQUEST,
				StatusMsg:  "评论文本为空",
			})
			return
		}
		comment := crud.CreateComment(db, &models.Comment{
			UserID:  user.ID,
			VideoID: request.VideoID,
			Content: request.CommentText,
		})
		ctx.JSON(http.StatusOK, response.CommentResponse{
			Response: response.Response{
				StatusCode: response.SUCCESS,
				StatusMsg:  "添加评论成功",
			},
			Comment: common.CommentModelChange(db, *comment),
		})
	} else {
		if request.CommentID == 0 {
			ctx.JSON(http.StatusBadRequest, response.Response{
				StatusCode: response.BADREQUEST,
				StatusMsg:  "删除失败",
			})
			return
		}
		crud.DeleteComment(db, request.CommentID)
		ctx.JSON(http.StatusOK, response.Response{
			StatusCode: response.SUCCESS,
			StatusMsg:  "评论删除成功",
		})
	}

}

func CommentList(ctx *gin.Context) {
	db := sql.GetSession()

	videoID := common.QueryIDAndValid(ctx, "video_id")
	if videoID == 0 {
		return
	}

	ctx.JSON(http.StatusOK, response.CommentListResponse{
		Response: response.Response{
			StatusCode: response.SUCCESS,
			StatusMsg:  "获取成功",
		},
		CommentList: common.CommentsModelChange(db, crud.GetComments(db, videoID)),
	})
}
